# Tmpl configuration

**Title:** Tmpl configuration

|                           |                                                           |
| ------------------------- | --------------------------------------------------------- |
| **Type**                  | `object`                                                  |
| **Required**              | No                                                        |
| **Additional properties** | [\[Not allowed\]](# "Additional Properties not allowed.") |

**Description:** A configuration file describing how a tmux session should be created.

| Property                        | Pattern | Type            | Deprecated | Definition               | Title/Description                                                      |
| ------------------------------- | ------- | --------------- | ---------- | ------------------------ | ---------------------------------------------------------------------- |
| - [tmux](#tmux)                 | No      | string          | No         | -                        | tmux executable                                                        |
| - [tmux_options](#tmux_options) | No      | array of string | No         | -                        | tmux command line options                                              |
| - [session](#session)           | No      | object          | No         | In #/$defs/SessionConfig | Session configuration describing how a tmux session should be created. |

## 1. Property `Tmpl configuration > tmux`

**Title:** tmux executable

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `"tmux"` |

**Description:** The tmux executable to use. Must be an absolute path, or available in $PATH.

## 2. Property `Tmpl configuration > tmux_options`

**Title:** tmux command line options

|              |                   |
| ------------ | ----------------- |
| **Type**     | `array of string` |
| **Required** | No                |

**Description:** Additional tmux command line options to add to all tmux command invocations.

**Examples:**

```yaml
['-f', '/path/to/tmux.conf']
```

```yaml
['-L', 'MySocket']
```

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be           | Description                                                                           |
| ----------------------------------------- | ------------------------------------------------------------------------------------- |
| [tmux_options items](#tmux_options_items) | A tmux command line flag and its value, if any. See 'man tmux' for available options. |

### 2.1. Tmpl configuration > tmux_options > tmux_options items

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

**Description:** A tmux command line flag and its value, if any. See `man tmux` for available options.

## 3. Property `Tmpl configuration > session`

|                           |                                                           |
| ------------------------- | --------------------------------------------------------- |
| **Type**                  | `object`                                                  |
| **Required**              | No                                                        |
| **Additional properties** | [\[Not allowed\]](# "Additional Properties not allowed.") |
| **Defined in**            | #/$defs/SessionConfig                                     |

**Description:** Session configuration describing how a tmux session should be created.

| Property                          | Pattern | Type   | Deprecated | Definition         | Title/Description                                                                                                                                                                                                                                                                               |
| --------------------------------- | ------- | ------ | ---------- | ------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [name](#session_name)           | No      | string | No         | In #/$defs/name    | A name for the tmux session or window. Must only contain alphanumeric characters, underscores, dots, and dashes                                                                                                                                                                                 |
| - [path](#session_path)           | No      | string | No         | In #/$defs/path    | The directory path used as the working directory in a tmux session, window, or pane.The paths are passed down from session to window to pane and can be customized at any level. If a path begins with '~', it will be automatically expanded to the current user's home directory. |
| - [env](#session_env)             | No      | object | No         | In #/$defs/env     | A list of environment variables to set in a tmux session, window, or pane.These variables are passed down from the session to the window and can be customized at any level. Please note that variable names should consist of uppercase alphanumeric characters and underscores.   |
| - [on_window](#session_on_window) | No      | string | No         | In #/$defs/command | On-Window shell command                                                                                                                                                                                                                                                                         |
| - [on_pane](#session_on_pane)     | No      | string | No         | In #/$defs/command | On-Pane shell command                                                                                                                                                                                                                                                                           |
| - [on_any](#session_on_any)       | No      | string | No         | In #/$defs/command | On-Window/Pane shell command                                                                                                                                                                                                                                                                    |
| - [windows](#session_windows)     | No      | array  | No         | -                  | Window configurations                                                                                                                                                                                                                                                                           |

### 3.1. Property `Tmpl configuration > session > name`

|                |                                              |
| -------------- | -------------------------------------------- |
| **Type**       | `string`                                     |
| **Required**   | No                                           |
| **Default**    | `"The current working directory base name."` |
| **Defined in** | #/$defs/name                                 |

**Description:** A name for the tmux session or window. Must only contain alphanumeric characters, underscores, dots, and dashes

| Restrictions                      |                                                                         |
| --------------------------------- | ----------------------------------------------------------------------- |
| **Must match regular expression** | `^[\w._-]+$` [Test](https://regex101.com/?regex=%5E%5B%5Cw._-%5D%2B%24) |

### 3.2. Property `Tmpl configuration > session > path`

|                |                                    |
| -------------- | ---------------------------------- |
| **Type**       | `string`                           |
| **Required**   | No                                 |
| **Default**    | `"The current working directory."` |
| **Defined in** | #/$defs/path                       |

**Description:** The directory path used as the working directory in a tmux session, window, or pane.

The paths are passed down from session to window to pane and can be customized at any level. If a path begins with '~', it will be automatically expanded to the current user's home directory.

**Examples:**

```yaml
/path/to/project
```

```yaml
~/path/to/project
```

```yaml
relative/path/to/project
```

### 3.3. Property `Tmpl configuration > session > env`

|                           |                                                                                                                         |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------- |
| **Type**                  | `object`                                                                                                                |
| **Required**              | No                                                                                                                      |
| **Additional properties** | [\[Should-conform\]](#session_env_additionalProperties "Each additional property must conform to the following schema") |
| **Defined in**            | #/$defs/env                                                                                                             |

**Description:** A list of environment variables to set in a tmux session, window, or pane.

These variables are passed down from the session to the window and can be customized at any level. Please note that variable names should consist of uppercase alphanumeric characters and underscores.

**Example:**

```yaml
APP_ENV: development
DEBUG: true
HTTP_PORT: 8080

```

| Property                                | Pattern | Type                      | Deprecated | Definition | Title/Description |
| --------------------------------------- | ------- | ------------------------- | ---------- | ---------- | ----------------- |
| - [](#session_env_additionalProperties) | No      | string, number or boolean | No         | -          | -                 |

#### 3.3.1. Property `Tmpl configuration > session > env > additionalProperties`

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string, number or boolean` |
| **Required** | No                          |

### 3.4. Property `Tmpl configuration > session > on_window`

**Title:** On-Window shell command

|                |                 |
| -------------- | --------------- |
| **Type**       | `string`        |
| **Required**   | No              |
| **Defined in** | #/$defs/command |

**Description:** A shell command to run first in all created windows. This is intended for any kind of project setup that should be run before any other commands. The command is run using the `send-keys` tmux command.

| Restrictions   |     |
| -------------- | --- |
| **Min length** | 1   |

### 3.5. Property `Tmpl configuration > session > on_pane`

**Title:** On-Pane shell command

|                |                 |
| -------------- | --------------- |
| **Type**       | `string`        |
| **Required**   | No              |
| **Defined in** | #/$defs/command |

**Description:** A shell command to run first in all created panes. This is intended for any kind of project setup that should be run before any other commands. The command is run using the `send-keys` tmux command.

| Restrictions   |     |
| -------------- | --- |
| **Min length** | 1   |

### 3.6. Property `Tmpl configuration > session > on_any`

**Title:** On-Window/Pane shell command

|                |                 |
| -------------- | --------------- |
| **Type**       | `string`        |
| **Required**   | No              |
| **Defined in** | #/$defs/command |

**Description:** A shell command to run first in all created windows and panes. This is intended for any kind of project setup that should be run before any other commands. The command is run using the `send-keys` tmux command.

| Restrictions   |     |
| -------------- | --- |
| **Min length** | 1   |

### 3.7. Property `Tmpl configuration > session > windows`

**Title:** Window configurations

|              |         |
| ------------ | ------- |
| **Type**     | `array` |
| **Required** | No      |

**Description:** A list of tmux window configurations to create in the session. The first configuration will be used for the default window.

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be        | Description                                                          |
| -------------------------------------- | -------------------------------------------------------------------- |
| [WindowConfig](#session_windows_items) | Window configuration describing how a tmux window should be created. |

#### 3.7.1. Tmpl configuration > session > windows > WindowConfig

|                           |                                                                             |
| ------------------------- | --------------------------------------------------------------------------- |
| **Type**                  | `object`                                                                    |
| **Required**              | No                                                                          |
| **Additional properties** | [\[Any type: allowed\]](# "Additional Properties of any type are allowed.") |
| **Defined in**            | #/$defs/WindowConfig                                                        |

**Description:** Window configuration describing how a tmux window should be created.

| Property                                                              | Pattern | Type    | Deprecated | Definition          | Title/Description                                                                                                                                                                                                                                                                               |
| --------------------------------------------------------------------- | ------- | ------- | ---------- | ------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [name](#session_windows_items_name)                                 | No      | string  | No         | In #/$defs/name     | A name for the tmux session or window. Must only contain alphanumeric characters, underscores, dots, and dashes                                                                                                                                                                                 |
| - [path](#session_windows_items_path)                                 | No      | string  | No         | In #/$defs/path     | The directory path used as the working directory in a tmux session, window, or pane.The paths are passed down from session to window to pane and can be customized at any level. If a path begins with '~', it will be automatically expanded to the current user's home directory. |
| - [command](#session_windows_items_command)                           | No      | string  | No         | In #/$defs/command  | A shell command to run within a tmux window or pane.The 'send-keys' tmux command is used to simulate the key presses. This means it can be used even when connected to a remote system via SSH or a similar connection.                                                             |
| - [commands](#session_windows_items_commands)                         | No      | array   | No         | In #/$defs/commands | A list of shell commands to run within a tmux window or pane in the order they are listed.If a command is also specified in the 'command' property, it will be run first.                                                                                                           |
| - [env](#session_windows_items_env)                                   | No      | object  | No         | In #/$defs/env      | A list of environment variables to set in a tmux session, window, or pane.These variables are passed down from the session to the window and can be customized at any level. Please note that variable names should consist of uppercase alphanumeric characters and underscores.   |
| - [active](#session_windows_items_active)                             | No      | boolean | No         | In #/$defs/active   | Whether a tmux window or pane should be selected after session creation. The first window and pane will be selected by default.                                                                                                                                                                 |
| - [panes](#session_windows_items_panes)                               | No      | array   | No         | -                   | Pane configurations                                                                                                                                                                                                                                                                             |
| - [additionalProperties](#session_windows_items_additionalProperties) | No      | object  | No         | -                   | -                                                                                                                                                                                                                                                                                               |

##### 3.7.1.1. Property `Tmpl configuration > session > windows > Window configuration > name`

|                |                  |
| -------------- | ---------------- |
| **Type**       | `string`         |
| **Required**   | No               |
| **Default**    | `"tmux default"` |
| **Defined in** | #/$defs/name     |

**Description:** A name for the tmux session or window. Must only contain alphanumeric characters, underscores, dots, and dashes

| Restrictions                      |                                                                         |
| --------------------------------- | ----------------------------------------------------------------------- |
| **Must match regular expression** | `^[\w._-]+$` [Test](https://regex101.com/?regex=%5E%5B%5Cw._-%5D%2B%24) |

##### 3.7.1.2. Property `Tmpl configuration > session > windows > Window configuration > path`

|                |                       |
| -------------- | --------------------- |
| **Type**       | `string`              |
| **Required**   | No                    |
| **Default**    | `"The session path."` |
| **Defined in** | #/$defs/path          |

**Description:** The directory path used as the working directory in a tmux session, window, or pane.

The paths are passed down from session to window to pane and can be customized at any level. If a path begins with '~', it will be automatically expanded to the current user's home directory.

**Examples:**

```yaml
/path/to/project
```

```yaml
~/path/to/project
```

```yaml
relative/path/to/project
```

##### 3.7.1.3. Property `Tmpl configuration > session > windows > Window configuration > command`

|                |                 |
| -------------- | --------------- |
| **Type**       | `string`        |
| **Required**   | No              |
| **Defined in** | #/$defs/command |

**Description:** A shell command to run within a tmux window or pane.

The 'send-keys' tmux command is used to simulate the key presses. This means it can be used even when connected to a remote system via SSH or a similar connection.

| Restrictions   |     |
| -------------- | --- |
| **Min length** | 1   |

##### 3.7.1.4. Property `Tmpl configuration > session > windows > Window configuration > commands`

|                |                  |
| -------------- | ---------------- |
| **Type**       | `array`          |
| **Required**   | No               |
| **Defined in** | #/$defs/commands |

**Description:** A list of shell commands to run within a tmux window or pane in the order they are listed.

If a command is also specified in the 'command' property, it will be run first.

**Example:**

```yaml
['ssh user@host', 'cd /var/logs', 'tail -f app.log']
```

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                  | Description                                              |
| ------------------------------------------------ | -------------------------------------------------------- |
| [command](#session_windows_items_commands_items) | A shell command to run within a tmux window or pane. ... |

##### 3.7.1.4.1. Tmpl configuration > session > windows > Window configuration > commands > command

|                |                 |
| -------------- | --------------- |
| **Type**       | `string`        |
| **Required**   | No              |
| **Defined in** | #/$defs/command |

**Description:** A shell command to run within a tmux window or pane.

The 'send-keys' tmux command is used to simulate the key presses. This means it can be used even when connected to a remote system via SSH or a similar connection.

| Restrictions   |     |
| -------------- | --- |
| **Min length** | 1   |

##### 3.7.1.5. Property `Tmpl configuration > session > windows > Window configuration > env`

|                           |                                                                                                                                       |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| **Type**                  | `object`                                                                                                                              |
| **Required**              | No                                                                                                                                    |
| **Additional properties** | [\[Should-conform\]](#session_windows_items_env_additionalProperties "Each additional property must conform to the following schema") |
| **Default**               | `"The session env."`                                                                                                                  |
| **Defined in**            | #/$defs/env                                                                                                                           |

**Description:** A list of environment variables to set in a tmux session, window, or pane.

These variables are passed down from the session to the window and can be customized at any level. Please note that variable names should consist of uppercase alphanumeric characters and underscores.

**Example:**

```yaml
APP_ENV: development
DEBUG: true
HTTP_PORT: 8080

```

| Property                                              | Pattern | Type                      | Deprecated | Definition | Title/Description |
| ----------------------------------------------------- | ------- | ------------------------- | ---------- | ---------- | ----------------- |
| - [](#session_windows_items_env_additionalProperties) | No      | string, number or boolean | No         | -          | -                 |

##### 3.7.1.5.1. Property `Tmpl configuration > session > windows > Window configuration > env > additionalProperties`

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string, number or boolean` |
| **Required** | No                          |

##### 3.7.1.6. Property `Tmpl configuration > session > windows > Window configuration > active`

|                |                |
| -------------- | -------------- |
| **Type**       | `boolean`      |
| **Required**   | No             |
| **Default**    | `false`        |
| **Defined in** | #/$defs/active |

**Description:** Whether a tmux window or pane should be selected after session creation. The first window and pane will be selected by default.

##### 3.7.1.7. Property `Tmpl configuration > session > windows > Window configuration > panes`

**Title:** Pane configurations

|              |         |
| ------------ | ------- |
| **Type**     | `array` |
| **Required** | No      |

**Description:** A list of tmux pane configurations to create in the window.

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                  | Description                                                      |
| ------------------------------------------------ | ---------------------------------------------------------------- |
| [PaneConfig](#session_windows_items_panes_items) | Pane configuration describing how a tmux pane should be created. |

##### 3.7.1.7.1. Tmpl configuration > session > windows > Window configuration > panes > PaneConfig

|                           |                                                           |
| ------------------------- | --------------------------------------------------------- |
| **Type**                  | `object`                                                  |
| **Required**              | No                                                        |
| **Additional properties** | [\[Not allowed\]](# "Additional Properties not allowed.") |
| **Defined in**            | #/$defs/PaneConfig                                        |

**Description:** Pane configuration describing how a tmux pane should be created.

| Property                                                      | Pattern | Type    | Deprecated | Definition          | Title/Description                                                                                                                                                                                                                                                                               |
| ------------------------------------------------------------- | ------- | ------- | ---------- | ------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [path](#session_windows_items_panes_items_path)             | No      | string  | No         | In #/$defs/path     | The directory path used as the working directory in a tmux session, window, or pane.The paths are passed down from session to window to pane and can be customized at any level. If a path begins with '~', it will be automatically expanded to the current user's home directory. |
| - [command](#session_windows_items_panes_items_command)       | No      | string  | No         | In #/$defs/command  | A shell command to run within a tmux window or pane.The 'send-keys' tmux command is used to simulate the key presses. This means it can be used even when connected to a remote system via SSH or a similar connection.                                                             |
| - [commands](#session_windows_items_panes_items_commands)     | No      | array   | No         | In #/$defs/commands | A list of shell commands to run within a tmux window or pane in the order they are listed.If a command is also specified in the 'command' property, it will be run first.                                                                                                           |
| - [env](#session_windows_items_panes_items_env)               | No      | object  | No         | In #/$defs/env      | A list of environment variables to set in a tmux session, window, or pane.These variables are passed down from the session to the window and can be customized at any level. Please note that variable names should consist of uppercase alphanumeric characters and underscores.   |
| - [active](#session_windows_items_panes_items_active)         | No      | boolean | No         | In #/$defs/active   | Whether a tmux window or pane should be selected after session creation. The first window and pane will be selected by default.                                                                                                                                                                 |
| - [horizontal](#session_windows_items_panes_items_horizontal) | No      | boolean | No         | -                   | Horizontal split                                                                                                                                                                                                                                                                                |
| - [size](#session_windows_items_panes_items_size)             | No      | string  | No         | -                   | Size                                                                                                                                                                                                                                                                                            |
| - [panes](#session_windows_items_panes_items_panes)           | No      | array   | No         | -                   | Pane configurations                                                                                                                                                                                                                                                                             |

##### 3.7.1.7.1.1. Property `Tmpl configuration > session > windows > Window configuration > panes > Pane configuration > path`

|                |                      |
| -------------- | -------------------- |
| **Type**       | `string`             |
| **Required**   | No                   |
| **Default**    | `"The window path."` |
| **Defined in** | #/$defs/path         |

**Description:** The directory path used as the working directory in a tmux session, window, or pane.

The paths are passed down from session to window to pane and can be customized at any level. If a path begins with '~', it will be automatically expanded to the current user's home directory.

**Examples:**

```yaml
/path/to/project
```

```yaml
~/path/to/project
```

```yaml
relative/path/to/project
```

##### 3.7.1.7.1.2. Property `Tmpl configuration > session > windows > Window configuration > panes > Pane configuration > command`

|                |                 |
| -------------- | --------------- |
| **Type**       | `string`        |
| **Required**   | No              |
| **Defined in** | #/$defs/command |

**Description:** A shell command to run within a tmux window or pane.

The 'send-keys' tmux command is used to simulate the key presses. This means it can be used even when connected to a remote system via SSH or a similar connection.

| Restrictions   |     |
| -------------- | --- |
| **Min length** | 1   |

##### 3.7.1.7.1.3. Property `Tmpl configuration > session > windows > Window configuration > panes > Pane configuration > commands`

|                |                  |
| -------------- | ---------------- |
| **Type**       | `array`          |
| **Required**   | No               |
| **Defined in** | #/$defs/commands |

**Description:** A list of shell commands to run within a tmux window or pane in the order they are listed.

If a command is also specified in the 'command' property, it will be run first.

**Example:**

```yaml
['ssh user@host', 'cd /var/logs', 'tail -f app.log']
```

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                              | Description                                              |
| ------------------------------------------------------------ | -------------------------------------------------------- |
| [command](#session_windows_items_panes_items_commands_items) | A shell command to run within a tmux window or pane. ... |

##### 3.7.1.7.1.3.1. Tmpl configuration > session > windows > Window configuration > panes > Pane configuration > commands > command

|                |                 |
| -------------- | --------------- |
| **Type**       | `string`        |
| **Required**   | No              |
| **Defined in** | #/$defs/command |

**Description:** A shell command to run within a tmux window or pane.

The 'send-keys' tmux command is used to simulate the key presses. This means it can be used even when connected to a remote system via SSH or a similar connection.

| Restrictions   |     |
| -------------- | --- |
| **Min length** | 1   |

##### 3.7.1.7.1.4. Property `Tmpl configuration > session > windows > Window configuration > panes > Pane configuration > env`

|                           |                                                                                                                                                   |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Type**                  | `object`                                                                                                                                          |
| **Required**              | No                                                                                                                                                |
| **Additional properties** | [\[Should-conform\]](#session_windows_items_panes_items_env_additionalProperties "Each additional property must conform to the following schema") |
| **Default**               | `"The window env."`                                                                                                                               |
| **Defined in**            | #/$defs/env                                                                                                                                       |

**Description:** A list of environment variables to set in a tmux session, window, or pane.

These variables are passed down from the session to the window and can be customized at any level. Please note that variable names should consist of uppercase alphanumeric characters and underscores.

**Example:**

```yaml
APP_ENV: development
DEBUG: true
HTTP_PORT: 8080

```

| Property                                                          | Pattern | Type                      | Deprecated | Definition | Title/Description |
| ----------------------------------------------------------------- | ------- | ------------------------- | ---------- | ---------- | ----------------- |
| - [](#session_windows_items_panes_items_env_additionalProperties) | No      | string, number or boolean | No         | -          | -                 |

##### 3.7.1.7.1.4.1. Property `Tmpl configuration > session > windows > Window configuration > panes > Pane configuration > env > additionalProperties`

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string, number or boolean` |
| **Required** | No                          |

##### 3.7.1.7.1.5. Property `Tmpl configuration > session > windows > Window configuration > panes > Pane configuration > active`

|                |                |
| -------------- | -------------- |
| **Type**       | `boolean`      |
| **Required**   | No             |
| **Default**    | `false`        |
| **Defined in** | #/$defs/active |

**Description:** Whether a tmux window or pane should be selected after session creation. The first window and pane will be selected by default.

##### 3.7.1.7.1.6. Property `Tmpl configuration > session > windows > Window configuration > panes > Pane configuration > horizontal`

**Title:** Horizontal split

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | No        |
| **Default**  | `false`   |

**Description:** Whether to split the window horizontally. If false, the window will be split vertically.

##### 3.7.1.7.1.7. Property `Tmpl configuration > session > windows > Window configuration > panes > Pane configuration > size`

**Title:** Size

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

**Description:** The size of the pane in lines for horizontal panes, or columns for vertical panes. The size can also be specified as a percentage of the available space.

**Examples:**

```yaml
20%
```

```yaml
50
```

```yaml
215
```

##### 3.7.1.7.1.8. Property `Tmpl configuration > session > windows > Window configuration > panes > Pane configuration > panes`

**Title:** Pane configurations

|              |         |
| ------------ | ------- |
| **Type**     | `array` |
| **Required** | No      |

**Description:** A list of tmux pane configurations to create in the pane.

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                              | Description                                                      |
| ------------------------------------------------------------ | ---------------------------------------------------------------- |
| [PaneConfig](#session_windows_items_panes_items_panes_items) | Pane configuration describing how a tmux pane should be created. |

##### 3.7.1.7.1.8.1. Tmpl configuration > session > windows > Window configuration > panes > Pane configuration > panes > PaneConfig

|                           |                                                           |
| ------------------------- | --------------------------------------------------------- |
| **Type**                  | `object`                                                  |
| **Required**              | No                                                        |
| **Additional properties** | [\[Not allowed\]](# "Additional Properties not allowed.") |
| **Same definition as**    | [Pane configuration](#session_windows_items_panes_items)  |

**Description:** Pane configuration describing how a tmux pane should be created.

##### 3.7.1.8. Property `Tmpl configuration > session > windows > Window configuration > additionalProperties`

|                           |                                                                             |
| ------------------------- | --------------------------------------------------------------------------- |
| **Type**                  | `object`                                                                    |
| **Required**              | No                                                                          |
| **Additional properties** | [\[Any type: allowed\]](# "Additional Properties of any type are allowed.") |

______________________________________________________________________

Generated using [json-schema-for-humans](https://github.com/coveooss/json-schema-for-humans)
