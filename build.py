#!/usr/bin/env python
import os
import sys, subprocess


def run_cmd(cmd):
    p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    stdout = p.communicate()[0].decode('utf-8').strip()
    return stdout


# Get last tag.
def get_version():
    return run_cmd('git describe --abbrev=0 --tags')


# Assemble build command.
def build_cmd():
    build_flag = []

    version = get_version()
    if version != "":
        build_flag.append("-X main.version '{}'".format(version))

    return 'go build -ldflags "{}"'.format(" ".join(build_flag))


def main():
    goos = "linux"
    if sys.argv[1] == "build":
        for arch in ["amd64", "arm64", "arm"]:
            cmd = 'GOARCH={arch} go build --ldflags "-X main.version={version}" -o mdns_{goos}_{arch} .'.format(
                arch=arch,
                version=get_version(),
                goos=goos
            )
            print(cmd)
            run_cmd(cmd)

    elif sys.argv[1] == "package":
        for arch in ["amd64", "arm64", "arm"]:
            cmd = 'tar -zcvf mdns_{os}_{arch}.tar.gz mdns_{os}_{arch} config.sample.d mdns.service'.format(
                arch=arch,
                os=goos,
            )
            print(cmd)
            run_cmd(cmd)


if __name__ == '__main__':
    main()
