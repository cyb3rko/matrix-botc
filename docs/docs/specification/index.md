---
comments: true
---

# [Draft] MBCS - Specification Version 1.0

This is the unofficial standard on how to define and process commands of bots for the [Matrix Protocol](https://matrix.org).

To better explain the details, we use the following example statement from [Draupnir](https://github.com/the-draupnir-project/Draupnir) as reference:  

Definition:
```
!dp list create <shortcode> <alias localpart>
```

Usage:
```
!dp list create mylist local-alias
```

## Components

In the following, the different components of a statement are listed and explained.  
All components should be separated from each other by a single whitespace.

### Prefix

The prefix is the first part of the statement, which identifies the recipient.  
This prefix helps a bot identify which statements he has to respond to and which to ignore.

It shall always have a leading exclamation mark, followed by a minimum of two and a maximum of ten letters and digits.  
A prefix shall be considered case-**in**sensitive.

**RegEx**: [`^![a-z0-9]{2,10}$`](https://regex101.com/r/EENJNk/1)  
**Example**:

- `!dp` (see [example](#))
- `!botv3`
- `!longprefix`

### Command

The command is the part of the statement that indicates what action to execute.  
Statements may contain one or more commands, located right after the prefix.  
_For the exception of a statement without a command, see [](#)_

It shall consist of a minimum of one and a maximum of twenty letters and digits.  
A command shall be considered case-**in**sensitive.  

(In the case of a statement without any command, so only the prefix itself, the command `help` shall be implicitly assumed.)

**RegEx**: [`^[a-z0-9]{1,20}$`](https://regex101.com/r/GED1SZ/1)  
**Examples**:

- `list` (see [example](#))
- `create` (see [example](#))
- `find2words`

### Command Chain

A command chain is the combination of multiple commands in a single statement.  
They should be separated by a single whitespace.

**Example**:

- `list create` (see [example](#))
- `config add`
- `status update`

### Arguments

A command or command chain may have one or more parameters, here called arguments.  
They shall be considered case-sensitive.  
This standard does not define valid characters for arguments, as you may pass anything that can be processed by the Matrix protocol, the homeserver and the bot logic.

Multiple arguments shall be separated by a single whitespace. If an argument contains one or more whitespaces, surround the argument with single quotation marks.

**Examples**:

- `mylist local-alias` (see [example](#))
- `!GKPWoymMiVrWlSLhud:matrix.org`
- `'example argument'`

## Processing

### Whitespace Handling

Before every processing step of a statement, leading and trailing whitespace shall be removed (trimmed).  
_(Whitespace here visualized as \_)_

1. `__!dp__list___create_mylist__local-alias_` shall be processed as  
`!dp__list___create_mylist__local-alias`
2. `__list___create_mylist__local-alias` shall be processed as  
`list___create_mylist__local-alias`
3. `___create_mylist__local-alias` shall be processed as  
`create_mylist__local-alias`
4. `_mylist__local-alias` shall be processed as  
`mylist__local-alias`

### Case Handling

Before every processing step and after the [Whitespace Handling](#whitespace-handling), the input shall be converted to lowercase, unless processing arguments (the only case-sensitive statement component).

That way the processor does not have to distinguish between e.g. the prefixes `!dp`, `!Dp`, `!dP` and `!DP`.

### Statement Recognition

A statement received by the bot shall only be processed if the trimmed statement...

- starts with the prefix and at least one whitespace,
- consists of only the prefix.

In all other cases, the message shall not be further processed.

### 
