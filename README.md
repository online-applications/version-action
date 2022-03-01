# Version Docker Action

This action calculates the next semVer version and outputs it.
The action searches for the following keywords, and bumps the version based on them:
```
breaking: major
feature: minor
bugfix: patch
```

## Outputs:
Next SemVer version.


**Required** Pull Request Env variables: The following env variables MUST be defined in your workflow
```
env:
  COMMIT_MESSAGE:         ${{ github.event.pull_request.body }}
  ENVIRONMENT:            ${{ github.event.pull_request.default_branch }}

```
**Required** Push Env variables: The following env variables MUST be defined in your workflow
```
env:
  COMMIT_MESSAGE:         ${{ github.event.head_commit.message }}
  ENVIRONMENT:            ${{ github.ref_name }}
```


## Example usage 1 - send a pre defined template of started message
**A previous checkout step with fetch-depth of 0 is required!!!**

```
    - name: Check out code
      uses: actions/checkout@v2
      with:
          fetch-depth: 0
    - name: Calculate version
        uses: online-applications/version-action@v1
        id: version
        run: echo Calculating version...
    - name: Use the new version to do something
        run: |
        git tag ${{ steps.version.outputs}
        git push --tags
```