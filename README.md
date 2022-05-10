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
### URLs

Formidable uses URLs to define how to handle schemas/defaults/values.

For example, to edit with Firefox a schema (in YAML) from an HTTPS server, while readig default values from `stdin` (in JSON) and using effective values from the local file system (in HCL), outputing updates to `stdout`:

```bash
echo '{}' | FORMIDABLE_BROWSER="firefox" frmd \
    edit
    --schema 'https://example.com/my-schema.yml' \
    --defaults 'stdin://local?format=json' \
    --values 'file:///my/file/absolute/path.hcl' \
    --output 'stdout://local?format=json'
```

### Available loaders

#### `stdin://`

> TODO: Write doc + example

#### `http://` and `https://`

> TODO: Write doc + example

#### `file://`

> TODO: Write doc + example

### Available formats

#### JSON

- **URL Query:** `?format=json`
- **File extension:** `.json`
- **As input:** yes
- **As output:** yes

#### YAML

- **URL Query:** `?format=yaml`
- **File extension:** `.yaml` or `.yml`
- **As input:** yes
- **As output:** yes

#### HCL

- **URL Query:** `?format=hcl`
- **File extension:** `.hcl`
- **As input:** yes
- **As output:** no

### Available outputs

#### `stdout://` (default)

> TODO: Write doc + example

#### `file://`

> TODO: Write doc + example

#### `exec://`

> TODO: Write doc + example

## Licence

AGPL-3.0
