## Comparing Dev Files and IDP

A non-exhaustive list of apparent differences between dev files and [IDP YAML schema](https://github.com/redhat-developer/odo-fork/blob/kdo-poc/docs/IDP-YAML.md). 

### Resources

#### Approximately Equivalent resource types between dev files and IDP YAML:
|Che Dev File|IDP Yaml|
|-|-|
|Command|Scenario|
|Action|Task|
|Component of type `dockerimage`|Container|
| N/A | Runtime|



#### Commands vs Scenarios:
- The closest dev file equivalent to our Scenarios are Commands. 
- IDP scenarios have fixed a `name` field that UDO will recognize from a fixed list ('full-build', 'incremental-build') and thus these names have semantic meaning; commands have a generic user-definable name field with no semantic meaning (they are for UI display only).
- UDO IDP ensures that only one scenario (eg `full-maven-build`) is running at a time, per project; commands have no such restrictions (FWICT), they are 'fire and forget'.
- Dev files/dev file commands are fixed for a particular workspace; in order to add new commands, you need to recreate the workspace (stopping and restarting is not enough). So, for example, if you have a Maven-focused workspace, you cannot use it to also run Node-focused dev file commands. In contrast, IDPs work at the project level, and you can run multiple simultaneous projects of different types (Java, Node, etc.)
- Commands are initiated directly by the user from within Cheia (they are exposed in the UI as Theia tasks), there is no automated watch-based build trigger. 


#### Actions vs Tasks:
- The closest dev file equivalent to our Tasks are Actions.
- The documentation claims that you can [only have 1 action per command](https://redhat-developer.github.io/devfile/devfile): `Now the only one command must be specified in list but there are plans to implement supporting multiple actions commands`. I didn't see any devfiles with commands that had multiple actions, so this may still be true. 
- I believe this necessarily means that only the equivalent of Runtime tasks are supported by dev files.
	- While it's true you can target any defined container with an action, in practice (since one command will only trigger one action) you must target the runtime container in order to have an affect on the application contents running in the container.
- Tasks can be shared between multiple scenarios; Actions cannot be shared between Commands (eg if you want the same action to run in multiple commands, you must define that action multiple times; OTOH since it's only 1 action per command, the impact is limited.)
- IDP YAML has `sourceMappings` that allows you to customize where the source is synchronized into the container. With dev files `the source is mounted on a location stored in the CHE_PROJECTS_ROOT environment variable that is made available in the running container of the image. This location defaults to '/projects'.`
- Che actions can modify source code; IDP tasks cannot modify source code (due to the use of a one-way sync).
- IDP tasks have a type field, which corresponds to where the task will run. There are 3 types of tasks (specified under the `type` field of `.spec.tasks`:
	- *Shared*: Tasks that run within a task container, and share that task container with another defined task.
	- *Standalone*: Tasks that run within a task container, but do not share that container wiht another defined task.
	- *Runtime*: Tasks that run within the runtime container; the runtime container is a special container reserved for the user's application. Runtime tasks may be used to, for example, configure the container, or start/stop/restart it. See more below.


#### Components vs Containers:
- The equivalent to our containers are dev file Components of the `dockerimage` type.
- IDP YAML has `idpRepoMappings` to allow you to customize an existing runtime/container image (eg add files to the standard maven image, or to a standard liberty image), thus the IDP developer is not required to ship their own custom container images (eg 'codewind-idp-maven'). There is no corresponding functionality in dev files (hence why, for example, the Maven image used by the Che Maven stack is a custom image from `quay.io/eclipse/che-java11-maven:nightly`, rather than the official `maven:3.6` images)
- You CAN override the entrypoint of a container via dev file, just like with IDP yaml. (`command: ['sleep', 'infinity']` in the component)
- Due to the limit of 1 action per command limit (see 'Actions vs Tasks' above), it appears that Actions are equivalent to our Runtime tasks (eg with the runtime defined as a container in a component).
	- If this restriction were to be eliminated, then we could likely simulate runtime/shared/standalone tasks with actions, communicating via shared volumes.
- No equivalent in dev files to the IDP YAML's differentiation between a runtime container and task container.
- Devfile volumes have no subpath field, which in the IDP case is used when volumes are mounted into containers.
- Dev file containers may contain [long running processes such as databases](https://github.com/eclipse/che-devfile-registry/blob/master/devfiles/java-mongo/devfile.yaml). At present, the IDP yaml does not define a way for long-lived processes such as database to run within independent containers (this is left to the user). 

#### Runtime:
- An IDP runtime is a container reserved for running the user's application: this may be a standalone application (Go, Node) or a runtime (Wildfly, Open Liberty). The runtime container is designed to provide a cleaner separation between the build environment (for example maven container images) and runtime environment (eg Wildfly production image).
- No direct equivalent concept to an IDP runtime in the dev file; the closest equivalent is to launch a container, based on a runtime image.
- Task container vs Runtime container:
  	- A task is a short-lived action (such as a build or a lint), which the UDO CLI will synchronously block on while it is running. Eg. If I have a task that calls `mvn package`, it will run for as long as the `mvn package` command executes within the container. 
		- Task containers are NOT disposed of once the short-lived action completes. We keep tasks containers around after the task completes, and reuse them for subsequent task invocations (so that we don't need to wait for a pod to mount volumes each time a task is running, though we like the idea of disposing of task containers after some amount of idle time)
	- A runtime is a long-lived action (an application or server runtime running in the container), which the UDO CLI will NOT synchronously block on while it is running.
	- Task containers are defined under `.spec.shared.containers`.
	- Runtime containers are defined under `.spec.runtime`.


#### Miscellaneous:
- No type detection in dev files (`typeDetection` section of IDP YAML)
- No file watching (`watched` section of IDP YAML)
- No upload filter (`uploadFilter` section of IDP YAML)
- Dev files have a number of pre-set environment variables, such as `CHE_PROJECTS_ROOT` (no IDP equivalent)
- Dev files include many fields which are unrelated to our UDO use case (eg the iterative building and deploying dev projects), including: projects to import into the workspace, chePlugins, vscode plugins, database containers to stand up, vscode launch actions (the vscode-task and vscode-launch types, vscode resources, other che editors, and more.


### Resources:

Official dev files:
https://github.com/eclipse/che-devfile-registry/tree/master/devfiles


'Introductory dev files' article in the official Che docs:
https://www.eclipse.org/che/docs/che-7/making-a-workspace-portable-using-a-devfile/

UML Diagram:
https://github.com/redhat-developer/devfile

Devfile schema:
https://redhat-developer.github.io/devfile/devfile
