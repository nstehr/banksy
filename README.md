# banksy

[![Build Status](https://travis-ci.org/nstehr/banksy.svg?branch=master)](https://travis-ci.org/nstehr/banksy) [![Coverage Status](https://coveralls.io/repos/github/nstehr/banksy/badge.svg?branch=master)](https://coveralls.io/github/nstehr/banksy?branch=master)



Simple application that will watch your github repos and label PRs, based on a simple rules file.

## Usage
1. Generate a token to be used when registering the banksy webhook
`./banksy -genToken`
2. Obtain a github API token for your repo.  This will be a personal access token, found under `Settings -> Developer Settings -> Personal Access Tokens`
3. Set the token from step 1 in the environment variable: `WEBHOOK_TOKEN`
4. Set your API token in the environment variable: `GITHUB_API_TOKEN`
5. define your rules file (more on that below).
6. Launch banksy
`./banksy -port <web hook server port> -baseURL <base url of github API> -rules <rules file path>`
`baseURL` is optional, and only used if you are using banksy with a github enterprise installation.
Banksy will listen for web hooks on `/banksy`.  
7. Go to the repo you would like banksy to watch and add a webook.  Don't forget to add the token you generated in step 1.

## Rules File
The rules that banksy uses to label PRs are of two types, a Glob Rule and a Size Rule.  These are defined in a yaml file 
and passed in as a command line parameter.  Example:

```yaml
rules: 
  - 
    type: "globRule"
    Globs: 
      - "*README*"
      - "*/foo/*"
    Label: "Foo"
  - 
    type: "globRule"
    Globs: 
      - "*/bar*"
    Label: "Bar"
  -
    type: "sizeRule"
    Compare: "greaterThan"
    NumFiles: 2
    NumChanges: 30
    Label: "Large"
  -
    type: "sizeRule"
    Compare: "lessThan"
    NumChanges: 5
    Label: "Small"
```

### Glob Rule
Uses a list of [globs](https://en.wikipedia.org/wiki/Glob_(programming)) to specify file paths in a PR that when
matched will generate a label to be added.  You can specify multiple globs to make up a rule, if any of the globs match
then the label will be applied.

```yaml
  - 
    type: "globRule"
    Globs: 
      - "*README*"
      - "*/foo/*"
    Label: "Foo"
```

### Size Rule
Will label based on the size of a PR.  Size can be calculated by either the number of changes in a PR, or the number of files
that make up the PR.  If both are specified, it will be calculated using an `or` of the two values.  

The size rule also takes a comparator, which can either be `lessThan` or `greaterThan`.  `lessThan` will check if the number of files
or number of changes is less than the specified value(s) and `greaterThan` will check if the number of changes is greater
than the specified values.

```yaml
  -
    type: "sizeRule"
    Compare: "greaterThan"
    NumFiles: 2
    NumChanges: 30
    Label: "Large"
  -
    type: "sizeRule"
    Compare: "lessThan"
    NumChanges: 5
    Label: "Small"
```
