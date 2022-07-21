# doc it?

At the moment this repo exists due to one reason - I need some practice to learn Golang.

So, task which I'm trying to solve here it's just find a way to provide generic documentation for 
any YAML files. There are really cool tool what I'm using at work [helm-docs](https://github.com/norwoodj/helm-docs).
But, as it's name states it works only with helm.

I'm writing a lot of YAML at work (A LOT). Usually it some configuration files (e.g. .gitlab-ci.yml).
I would like to document them directly (in comments) and render documentation based on those comments.

This:

```yaml
# @doc-it
# # Hello my friend
# List of usefully things:
# * first
# * second
# * third
stages:
  - build
  - test

# @doc-it
# Text text
# [link](google.com)
build-code-job:
  stage: build
  script:
    - echo "Check the ruby version, then build some Ruby project files:"
    - ruby -v
    - rake
```

Will be rendered to something like this:

```markdown
# Hello my friend
List of usefully things:
* first
* second
* third

Text text
[link](google.com)
```

TODO:

- [ ] provide reference to YAML node
- [ ] put documented content to foldable block under documentation
