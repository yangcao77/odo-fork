## IDP YAML (Draft)

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

    env: # Defined As below

    volumeMappings: #  Optional: ability to map paths in the container to persistent volume paths
    - volumeName: idp-data-volume
      containerPath: /some/path/idp-data

    kubernetes: # Values only used for Kube deployments
    
      # TODO: Are there other Kube resource parameters we need to include here? securityContext? (cluster) role bindings?
      livenessProbe: # Optional, otherwise sane defaults should be used
        initialDelaySeconds: 15
        timeoutSeconds: 60

      readinessProbe: # Optional, otherwise sane defaults should be used
        initialDelaySeconds: 15
        timeoutSeconds: 60

      requests: # Are these optional in Kube?
        memory: "64Mi"
        cpu: "250m"
      limits: # Are these optional in Kube?
        memory: "128Mi"
        cpu: "500m"

  shared:

    containers:
    - name: maven-build-container
      image: docker.io/maven:3.6
      
      volumeMappings: #  Optional: ability to map paths in the container to persistent volume paths
      - volumeName: idp-data-volume
        containerPath: /some/path/idp-data
      # Map a directory for the task to copy data to runtime, or for some other arbitrary purpose
  
      env: # Optional key/value env var pairs, as above

      kubernetes: # Optional
        # Defined same as above

    
    # tasks: removed on (09/20), as we have hardcoded specific defaults for these, at this time. See 'Tasks' below.
    
    # tasks: 

      # If true, tasks that share the same build image will NOT run within the same container during a scenario.
      # If false, tasks that share the same build image WILL run in the same container during a scenario.
      # Note: Whether the container will be disposed after the scenario has completed is determined by disposeOnScenarioComplete.
      # Optional: default is true.
      # disposeOfSharedContainersOnTaskComplete: true # (true/false) 
      # (09/20) Not in plan to support at this time.
      
      # Whether a task container will be disposed of after the scenario has completed.
      # If true, all containers that were used in a scenario will be destroyed once the scenario ends. 
      # If false, all containers that were used in a scenario will be preserved for the next run.
      # Note: this applies BOTH to tasks that share a build image with another task, and those that don't.
      # Optional: default is true.
      # disposeOnScenarioComplete: true # (true/false) 
      # (09/20) Not in plan to support at this time.

      # Number of seconds to keep a task container alive, if the task container is not invoked during that period.
      # Optional: default is no timeout.
      # idleTaskContainerTimeout: 3600 
      # (09/20) Not in plan to support at this time.
    
    volumes: 
    - name: idp-data-volume 
      # UDO will decide how to create the volume (RWO/RWX) based on how many tasks reference the container (if >1, then RWX)
      size: 1Gi # kube only

    env: # Optional: Ability to map key/value pairs into the container as environment variables, shared between both runtime and tasks
    - name: key
      value: value

  tasks:
    # Task containers will ALWAYS stay up and be reused after they are used (eg they will never be disposed of after a single use).

    # Tasks that share the same build image will ALWAYS run in the same container during a scenario.
    # - Tasks that share a build image value, must have the exact same volume mappings, or the IDP is invalid and should not be executed.

    - name: maven-build
      container: maven-build-container
      command: /scripts/build.sh # could also just be a normal command ala `mvn clean package`
      # Tasks containers will always be started with a command to tail -f /dev/null, so that they persist. The actual tasks themselves will be run w/ kubectl exec
      
      workingDirectory: /codewind-workspace-mount-point # optional, where in the container to run the command

      logs: # Ability to reference arbitrary log file types that aren't included in container stderr/stdout
      - type: maven.build
        path: /logs/(etc)
    
      repoMappings: # Optional: Automatically upload files/directories from the IDP repo to a container on/before startup
      - srcPath: "/resources/scripts/build.sh"
        destPath: "/scripts/build.sh"
        setExecuteBit: true # Set execute bit on a single file
      - srcPath: "/scripts2" # path in remote git repo, where folder containing "idp.yaml" is /
        destPath: "/home/user" # path inside container to upload the directory
        setExecuteBit: true # Set execute bit on all files in the directory (required for windows local, git repos without execute, http serving)

      sourceMapping: # Optional: Ability to map files in the local project directory (eg the user's current working dir)` into the container
        destPath: "/home/user/src" # path inside container to copy the folder
        setExecuteBit: true # Set execute bit on all files in the directory
      # This is used to know where the source files should be copied into the container, might be useful for other scenarios like customization
      # Path should be a valid path within the container (but if volumes are mapped into paths in the container, you can use those volume paths)
              
      env: # Optional key/value env var pairs, as above
      # Values specified here will replace those specified in container, if there is an overlap.

    - name: server-start
      command: /opt/ibm/wlp/bin/server start $SERVER 
      # Validation Constraint: If a task doesn't specifiy a build image field (indicating a task should run in the runtime container), then that task should not have volume mappings. 
      # The task will still be able to use whichever mappings are specified in the runtime volumeMapping/
      # For example, if this 'server-start' task had volume mappings, the UDO tool should fail to run it.
      
  scenarios:
    - name: full-build
      tasks: ["maven-build", "server-start"]
    - name: incremental-build
      tasks: ["incremental-maven-build", "server-start"] # incremental-maven-build not actually defined in this sample

```

#### Update History:
- September 16th: Remove `kind`, update `github-id` to `githubId`
- September 17th: Remove `buildImage: docker.io/maven:3.6` from `server-start`, remove `maven-cache-volume` volume, removed `spec.shared.volumes.labels` and `spec.shared.volumes.accessModes`. `env` updated to Kubernetes-style `key/value`, to allow easy parsing.
- September 20th: 
  - `sourceMappings` -> `sourceMapping`, and removed the `srcPath` field (will always sync from project root).
  - Removed `spec.shared.tasks`, and all the fields under it, as we have hardcoded defaults for these values.
  - Added ability to map volumes into runtime image (this was always implied, but is now included), under `spec.runtime.volumeMappings`
- September 23rd:
  - `runAsUser` removed from `spec.tasks`
  - `buildImage` renamed to `image` under `spec.tasks`
  - Replace previous memory limit with `kubernetes.requests` and `kubernetes.limits`
