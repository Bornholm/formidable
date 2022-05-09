# Formidable

Web-based - but terminal compatible ! - little form editor using [JSONSchema](https://json-schema.org/) models.

**Terminal compatible ?**

The generated web UI does not use Javascript and can be used with terminal browsers like [w3m](https://en.wikipedia.org/wiki/W3m) or [lynx](https://en.wikipedia.org/wiki/Lynx_(web_browser)).

## Install

### Manually

Download the pre-compiled binaries from the [releases page](https://github.com/Bornholm/formidable/releases) and copy them to the desired location.

### Bash script

```
curl -sfL https://raw.githubusercontent.com/Bornholm/formidable/master/misc/script/install.sh | bash
```

It will download `frmd` to your current directory.

#### Script available environment variables

|Name|Description|Default|
|----|-----------|-------|
|`FORMIDABLE_VERSION`|Formidable version to download|`latest`|
|`FORMIDABLE_DESTDIR`|Formidable destination directory|`.`|

## Usage

### Defining a schema

### Using default values

### Handling values update

## Licence

AGPL-3.0