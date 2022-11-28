import git
from pathlib import Path
import semver
from sys import argv

import subprocess

STARTING_SEMVER = semver.parse("0.3.10")

# If --rebuild-all is passed, we will rebuild all schema files, even if they
# already exist in the schema.bacalhau.org directory
rebuild_all = False

if len(argv) > 1 and argv[1] == "--rebuild-all":
    rebuild_all = True

rootPath = Path(__file__).parent.parent.parent

repo = git.Repo(rootPath)
repo.branches.main.checkout()

tagList = repo.git.ls_remote("--tags", "origin").split("refs/tags/")[1:]

listOfTagsToBuild = []

for longTag in tagList:
    splitValues = longTag.strip().split("\n")
    if len(splitValues) > 1:
        tag, commit = splitValues
    else:
        tag = splitValues[0]
        commit = None

    if tag.startswith("v"):
        tag = tag[1:]

    try:
        semVerTag = semver.VersionInfo.parse(tag)
        print(semVerTag)
        if semVerTag > STARTING_SEMVER:
            listOfTagsToBuild.append(tag)
    except ValueError as ve:
        print(f"Skipping {tag} because it is not a valid semver tag: {ve}")
        continue

for tag in listOfTagsToBuild:
    repo.active_branch.checkout(f"v{tag}")
    subprocess.call(["make", "release"], cwd=rootPath)
    output = subprocess.call(["bin/bacalhau", "validate", "--output-schema"], cwd=rootPath, stdout=subprocess.PIPE)
