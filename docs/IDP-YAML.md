## IDP YAML

```yaml
apiVersion: codewind.dev/v1alpha1
metadata: 
  name: java-liberty-pack
  language: java 
  framework: liberty 
  version: 0.0.1 
  maintainers:  
  - name: Jonathan West
    email: jgwest@email.com
    githubID: jgwest

spec:

  dev: 
    watched:
      # On the local machine, which path to watch, from the root folder of the project.
      # For example /codewind-workspace/my-project would be /
      watchDir: /src # Optional, otherwise assumed to be the root /
      ignoredPaths: # Same format as filewatcher, Optional, if not specified then no excludes.
      - /target
      - /some.file

    uploadFilter: # Unclear if we can just combine this with watch, or if it needs to be separate
      ignoredPaths: # Same format as filewatcher, Optional, if not specified then no excludes.
      - /target

    typeDetection: # Optional: Rules that are used to detect if this IDP is applicable to a given project (eg OpenLiberty IDP for an OpenLiberty-based app)
    # At least one rule must match for a type; ie this is an OR clause list, not an AND clause list.
    - pathSelector: # At least one 'pathSelector' entry must exist if typeDetection is specified
        # Optional, Either a 'rootPath' xor a 'filenameWildcard' must be specified
        rootPath: # example: /pom.xml, or /go.mod
        filenameWildcard: # example: server.xml, or *.go; this means recursively walk a project and look for filenames that match this pattern. Same pattern style as filewatcher.
      textStringsToMatch: # Optional: If a file matches one of the selectors, then look for at least one of these strings (ie OR clause list, not AND clause list). 
      - net.wasdev.wlp.maven.plugins # Example: if this ID is found in the file, then this IDP should be considered to be applicable to the project
    
  runtime:
    
    image: docker.io/ibmcom/websphere-liberty:latest

    endpoints: # Not optional if HTTP(S) port is specified
      containerHealth: /health # How to tell the container is healthy
      appRoot: /app # Not a health check

    ports: # All are optional, display a warning if neither HTTP port is specified
      internalHttpPort: 9080 
      internalHttpsPort: 9443
      internalDebugPort: 7777
      internalPerformancePort: 9050
      
    logs: # Optional: Ability to reference arbitrary log file types that aren't included in container stderr/stdout
    - type: some.type
      path: /logs/(etc)

    env: # As below

    kubernetes: # Values only used for Kube deployments
    
      # TODO: Are there other Kube resource parameters we need to include here? securityContext? (cluster) role bindings?
      livenessProbe: # Optional, otherwise sane defaults should be used
        initialDelaySeconds: 15
        timeoutSeconds: 60

      readinessProbe: # Optional, otherwise sane defaults should be used
        initialDelaySeconds: 15
        timeoutSeconds: 60

      memoryLimit: 600Mi # Arguable whether this should only be in kubernetes 

    # runAsUser: 185 

  shared:
    tasks: # Optional, see defaults below

      # If true, tasks that share the same build image will NOT run within the same container during a scenario.
      # If false, tasks that share the same build image WILL run in the same container during a scenario.
      # Note: Whether the container will be disposed after the scenario has completed is determined by disposeOnScenarioComplete.
      # Optional: default is true.
      disposeOfSharedContainersOnTaskComplete: true # (true/false) 

      # Whether a task container will be disposed of after the scenario has completed.
      # If true, all containers that were used in a scenario will be destroyed once the scenario ends. 
      # If false, all containers that were used in a scenario will be preserved for the next run.
      # Note: this applies BOTH to tasks that share a build image with another task, and those that don't.
      # Optional: default is true.
      disposeOnScenarioComplete: true # (true/false) 

      # Number of seconds to keep a task container alive, if the task container is not invoked during that period.
      # Optional: default is no timeout.
      idleTaskContainerTimeout: 3600 
    
    volumes: 
    - name: idp-data-volume 
      # UDO will decide how to create the volume (RWO/RWX) based on how many tasks reference the container (if >1, then RWX)
      size: 1Gi # kube only

    env: # Optional: Ability to map key/value pairs into the container as environment variables, shared between both runtime and tasks
    - name: key
      value: value

  tasks:
    - name: maven-build
      buildImage: docker.io/maven:3.6
      command: /scripts2/build.sh # could also just be a normal command ala `mvn clean package`
      workingDirectory: /codewind-workspace-mount-point # optional, where in the container to run the command

      logs: # Ability to reference arbitrary log file types that aren't included in container stderr/stdout
      - type: maven.build
        path: /logs/(etc)
    
      volumeMappings: #  Optional: ability to map paths in the container to persistent volume paths
      - volumeName: idp-data-volume
        containerPath: /idp-data-voume
      # Map a directory for the build job to be able to copy the .war file

      repoMappings: # Optional: Automatically upload files/directories from the IDP repo to a container on/before startup
      - srcPath: "/scripts" # path in remote git repo, where folder containing "idp.yaml" is /
        destPath: "/home/user" # path inside container to upload the directory
        setExecuteBit: true # Set execute bit on all files in the directory (required for windows local, git repos without execute, http serving)
      - srcPath: "/scripts2/build.sh"
        destPath: "/home/user/build-scripts"
        setExecuteBit: true # Set execute bit on a single file

      sourceMappings: # Optional: Ability to map files in the local project directory (eg the user's current working dir)` into the container
      - srcPath: "/src" # copy from $CURRENT_DIR/src
        destPath: "/home/user/src" # path inside container to copy the folder
        setExecuteBit: true # Set execute bit on all files in the directory
      # This is used to know where the source files should be copied into the container, might be useful for other scenarios like customization
              
      env: # Optional key/value env var pairs, as above

      runAsUser: 185 # Optional, same as above
      
      kubernetes: # Optional, as above
        livenessProbe: 
        readinessProbe:
        memoryLimit: 600Mi # As above
      
    - name: server-start
      command: /opt/ibm/wlp/bin/server start $SERVER 

  scenarios:
    - name: full-build
      tasks: ["maven-build", "server-start"]
    - name: incremental-build
      tasks: ["incremental-maven-build", "server-start"]
```

#### Update History:
- September 16th: Remove `kind`, update `github-id` to `githubId`
- September 17th: Remove `buildImage: docker.io/maven:3.6` from `server-start`, remove `maven-cache-volume` volume, removed `spec.shared.volumes.labels` and `spec.shared.volumes.accessModes`. `env` updated to Kuberenetes-style `key/value`, to allow easy parsing.
