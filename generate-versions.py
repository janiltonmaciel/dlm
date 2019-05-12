# -*- encoding: utf-8 -*-
import codecs
import fnmatch
import os
import re
import sys
from datetime import datetime
from io import StringIO
from os.path import basename, dirname, exists, join
from urllib.parse import urlparse

import requests
from natsort import natsorted

from ruamel.yaml import YAML


class Hub():

    def __init__(self, language):
        self.language = language
        self.images = []
        self.all_images = []

    def add_image(self, image):
        self.all_images.append(image)

        # for data in self.images:
        #     inter = image.stags.intersection(data.stags)
        #     if len(inter) > 0:
        #         if image.git_datetime < data.git_datetime:
        #             for tag in inter:
        #                 image.stags.remove(tag)

        if len(image.stags) > 0:
            self.images.append(image)

    def print_all_versions(self):
        for image in self.all_images:
            print(image.version, image.tags, image.distribution, image.release)

    def print_versions(self):
        versions = self._versions()
        for v in versions:
            print(v.version)

    def save(self, filename=None):
        file_name = filename or "config/{}-versions.yml".format(self.language)
        print("Save {}...".format(file_name))

        data = []
        versions = self._versions()
        for v in versions:
            stags = set(v["version"].stags)
            for vs in v["distributions"]:
                stags.update(list(vs.stags))

            distributions_hash = {}
            distributions = []
            for vs in v["distributions"]:
                key = "{release_name}-{release}-{image}".format(release_name=vs.distribution["name"], release=vs.distribution["release"], image=vs.image_name)
                if not key in distributions_hash:
                    distributions.append({
                        "name": vs.distribution["name"],
                        "releaseName": vs.distribution["release_name"],
                        "release": float(vs.distribution["release"]),
                        "weight": vs.distribution["weight"],
                        "image": vs.image_name,
                        "tags": natsorted(list(vs.stags)),
                        "urlRepository": vs.url_repository,
                        "urlDockerfile": vs.url_dockerfile,
                    })
                    distributions_hash[key] = True

            data_version = {
                "version": v["version"].version,
                "majorVersion": v["version"].major_version,
                "prerelease": v["version"].prerelease,
                "date": v["version"].git_datetime.isoformat(),
                # "tags": natsorted(list(stags)),
                # "distributionReleases": ", ".join(natsorted(set(["{} {}".format(vs.distribution["name"], vs.distribution["release_name"]) for vs in v["distributions"]]))),
                "distributionReleases": ", ".join(natsorted(set([vs["name"].title() for vs in distributions]))),
                "distributions": distributions,
            }
            data.append(data_version)

        save_yml(file_name, data)

    def _versions(self):
        images_versions = {"{}${}".format(image.version, image.url_dockerfile):image for image in self.images}
        versions = list(images_versions.keys())

        group = {}
        versions = natsorted(versions, reverse=True)
        for version in versions:
            v, url_dockerfile = version.split("$")
            value = group.get(v)
            if not value:
                value = []
            value.append(version)
            group[v] = value

        # import pdb; pdb.set_trace()
        data = []
        existe = {}
        for version in versions:
            v, url_dockerfile = version.split("$")
            value = group[v]
            if v in existe:
                continue

            data.append({
                "version": images_versions[version],
                "distributions": [images_versions[x] for x in value],
            })
            existe[v] = True


        return data


class Distribution():
    DISTRIBUTIONS = {}

    DEBIAN = {
        "SQUEEZE": {"name": "Debian", "release_name": "6 (squeeze)", "release": 6},
        "WHEEZY": {"name": "Debian", "release_name": "7 (wheezy)", "release": 7},
        "JESSIE": {"name": "Debian", "release_name": "8 (jessie)", "release": 8},
        "STRETCH": {"name": "Debian", "release_name": "9 (stretch)", "release": 9},
        "BUSTER": {"name": "Debian", "release_name": "10 (buster)", "release": 10},
    }

    UBUNTU = {
        "TRUSTY": {"name": "Ubuntu", "release_name": "14.04 (trusty)", "release": 14.04},
        "XENIAL": {"name": "Ubuntu", "release_name": "16.04 (xenial)", "release": 16.04},
        "BIONIC": {"name": "Ubuntu", "release_name": "18.04 (bionic)", "release": 18.04},
        "COSMIC": {"name": "Ubuntu", "release_name": "18.10 (cosmic)", "release": 18.10},
        "DISCO": {"name": "Ubuntu", "release_name": "19.04 (disco)", "release": 19.04},
    }

    UBUNTU.update({
        "{}".format(d["release"]):d
        for _, d in UBUNTU.items()
    })

    ALPINE = [3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 3.8, 3.9, 4.0]

    # ADD DEBIAN / UBUNTU / BUILD_PACK_DEPS
    BUILD_PACK_DEPS = {
        **DEBIAN,
        **UBUNTU
    }
    DISTRIBUTIONS["buildpack-deps"] = {**DEBIAN["STRETCH"], "weight": 5}
    for k, d in BUILD_PACK_DEPS.items():
        key_buildpack ="buildpack-deps:{}".format(k).lower()
        key_buildpack_scm = "{}-scm".format(key_buildpack)
        key_buildpack_curl = "{}-curl".format(key_buildpack)
        key_debian = "{}:{}".format(d["name"], k).lower()
        key_debian_slim = "{}:{}-slim".format(d["name"], k).lower()

        # buildpack-deps:stretch > buildpack-deps:stretch-scm > buildpack-deps:stretch-curl > debian:stretch > debian:stretch-slim
        weight = 5
        for key in [key_buildpack, key_buildpack_scm, key_buildpack_curl, key_debian, key_debian_slim]:
            DISTRIBUTIONS[key] = {**d, "weight": weight}
            weight = weight - 1

    # ADD ALPINE
    DISTRIBUTIONS.update({
        "alpine:{}".format(d).lower():{"name": "Alpine", "release_name": str(d), "release": d, "weight": 1}
        for d in ALPINE
    })

    @classmethod
    def get_distributions(cls):
        return cls.DISTRIBUTIONS

    @classmethod
    def save(cls):
        file_name = "config/distributions.yml"
        print("Save {}...".format(file_name))

        data = {}
        for distribution in cls.get_distributions().values():
            key = "{name}-{release_name}-{release}".format(**distribution)
            if not key in data:
                data[key] = {
                    "name": distribution["name"],
                    "release_name": distribution["release_name"],
                    "release": distribution["release"],
                }
        data = sorted(data.values(), key=lambda x: x["name"] and x["release"], reverse=True)
        save_yml(file_name, data)


class Image():
    RE_PRERELEASE = re.compile(r"(rc|beta|a|b)+\d+$")
    RE_IGNORE = re.compile(r"(-preview|onbuild|windowsservercore|nanoserver|-cross|chakracore)")
    RE_FROM_IMAGE = re.compile(r"^(?!#)\s*FROM\s+([-\w.:/]+)", re.MULTILINE)
    RE_RELEASE = re.compile(r"([\.\d]+)$")
    RE_MAJOR_VERSION = re.compile(r'^(\d+)\.(\d+)')
    RE_VERSIONS = {
        # "ruby": re.compile(r"^(([.\d]+)(-p\d+)?)"),
        # "python": re.compile(r"^(([.\d]+)(\w\d+)?)"),
        # "golang": re.compile(r"^(([.\d]+)((beta|rc)\d+)?)"),
        # "default": re.compile(r"^([.\d]+)"),
        "default": re.compile(r"^(\d+[^-]+)"),
    }

    def __init__(self, language, git_datetime, tags, directory, git_commit, repository, filepath):
        self.language = language
        self.tags = tags
        self.stags = set(map(str.strip, tags.split(",")))
        self.git_datetime = git_datetime
        self.git_commit = git_commit
        self.directory = directory
        self.repository = repository
        self.filepath = filepath

        self.re_versions = self.RE_VERSIONS.get(self.language) or self.RE_VERSIONS["default"]
        self.url_repository = "https://github.com{}".format(repository)
        self.url_dockerfile = self.format_url()

        self.distribution = None
        self.version = None
        self.major_version = None
        self.image_name = None
        self.valid = True

        self.ignore = self.RE_IGNORE.search(self.tags) != None

        if not self.ignore:
            self._proccess()


    def _proccess(self):
        try:
            self.version, self.major_version = self._versions(self.tags)
        except Exception as e:
            print(e)

        if not self.version or not self.major_version:
            log_debug("INVALID VERSION:", date=self.git_datetime.isoformat(), tags=self.tags)
            self.valid = False
            return

        self.distribution, self.image_name = self.get_distribution_release()

        if not self.distribution or not self.image_name:
            print("INVALID IMAGE:", self.image_name, self.url_dockerfile)
            self.valid = False
            return

    def _versions(self, tags):
        stags = set(map(str.strip, tags.split(",")))
        versions = []
        for tag in stags:
            resp = self.re_versions.search(tag)
            if resp :
                versions.append(resp.groups()[0])

        if not versions:
            return None, None

        version = max(versions, key=len)

        match = self.RE_MAJOR_VERSION.match(version)
        if match:
            vmajor, vminor = match.group(1, 2)
            major = "{}.{}".format(vmajor, vminor)
        else:
            major = version[0:3]

        return version, major

    def get_distribution_release(self):
        text = self.download(self.url_dockerfile)
        # import pdb; pdb.set_trace()
        resp = self.RE_FROM_IMAGE.search(text)
        if resp:
            image_name = resp.groups()[0]
            # print("IMAGE:", image_name)
            # print("TAGS:", image_name, " - ", self.tags)
            # print(self.tags)
            return self.find_distribution_release(image_name)

        print("REGEX IMAGE NOT FOUND", self.url_dockerfile)
        return None, None

    def find_distribution_release(self, image_name):
        stags = set(map(str.strip, image_name.split(",")))
        for key, distribution in Distribution.get_distributions().items():
            if image_name == key:
                return distribution, image_name

        return None, None

    def is_valid(self):
        return not self.ignore and self.valid

    @property
    def prerelease(self):
        return self.RE_PRERELEASE.search(self.version) != None

    def request_url(self, url):
        with requests.get(url) as r:
            r.raise_for_status()
            return r.text

    def download(self, url):
        # import pdb; pdb.set_trace()
        tmp_filepath = "tmp/download/{}/{}".format(self.language, self.slugify(url))
        if not exists(tmp_filepath):
            try:
                self.download_file(url, tmp_filepath)
            except Exception as e:
                print("Error %s", str(e))
                print(self.__dict__)
                return ""


        return self.load(tmp_filepath)

    def download_file(self, url, file_path):
        dir_path = dirname(file_path)
        if not exists(dir_path):
            os.makedirs(dir_path)

        with requests.get(url, stream=True) as r:
            r.raise_for_status()
            with open(file_path, 'wb') as f:
                for chunk in r.iter_content(chunk_size=512):
                    if chunk:
                        f.write(chunk)

    def load(self, file_path):
        with codecs.open(file_path, 'r', encoding='utf8') as fp:
            return fp.read()

    def format_url(self):
        raw_git_url = "https://raw.githubusercontent.com"
        return "{raw_git_url}{repository}/{commit}/{directory}Dockerfile".format(
            raw_git_url=raw_git_url,
            repository=self.repository,
            commit=self.git_commit,
            directory= "%s/" % self.directory if self.directory else "")

    def slugify(self, s):
        s = s.replace('/', '-')
        s = re.sub(r'&[a-z]+;', '', s)
        s = re.sub(r'[^\$\w\s-]', '', s)
        s = re.sub(r'[-\s]+', '-', s)[:]
        return s.strip('-')



class OfficialImages():
    FIELDS = {
        "tags": {"pattern": r"Tags:(.*)", "format": str},
        "git_commit": {"pattern": r"^GitCommit:(.*)", "format": str},
        "directory": {"pattern": r"^Directory:(.*)", "format": str},
    }

    def __init__(self, dir_path, language):
        log_debug(dir_path, language)

        self.hub = Hub(language)
        self.language = language
        self.dir_path = dir_path
        self.files = self._files()
        self.data = self._load()
        self._parser()

    def print_versions(self):
        return self.hub.print_versions()

    def print_all_versions(self):
        return self.hub.print_all_versions()

    def save(self, filename=None):
        return self.hub.save(filename)

    def _parser(self):
        for d in self.data:
            repository = self.get_value(d["data"], r"^GitRepo:(.*)")
            blocks = self._blocks(content=d["data"], re_start=r"^Tags:", re_stop=r"^\s+")
            if blocks:
                for block in blocks:
                    image_kwargs = {
                        "language": self.language,
                        "git_datetime": d["git_datetime"],
                        "repository": self._repository(repository),
                        "filepath": d["filepath"],
                    }
                    fields_kwargs = self.get_value_fields(block, self.FIELDS, ["tags", "git_commit", "directory"])
                    image_kwargs.update(fields_kwargs)

                    image = Image(**image_kwargs)
                    if image.is_valid():
                        self.hub.add_image(image)
            else:
                for line in d["data"]:
                    resp = re.search(r"^([.\w-]+):\s+(.*)@([\w]+)+\s+(.*)$", line)
                    if resp:
                        tags, repository, git_commit, directory = resp.groups()

                        image_kwargs = {
                            "language": self.language,
                            "git_datetime": d["git_datetime"],
                            "tags": tags,
                            "git_commit": git_commit,
                            "directory": directory,
                            "repository": self._repository(repository),
                            "filepath": d["filepath"],
                        }
                        image = Image(**image_kwargs)
                        if image.is_valid():
                            self.hub.add_image(image)


    def _load(self):
        data = []
        for f in self.files:
            with codecs.open(f, 'r', encoding='utf8') as fp:
                data.append({
                    "git_datetime": self._git_datetime(f),
                    "data": fp.readlines(),
                    "filepath": f,
                    "filename": basename(f),
                })

        return data

    def _git_datetime(self, filepath):
        filename = basename(filepath)
        dt = filename.split("_")[1]
        return datetime.strptime(dt, "%Y-%m-%dT%H:%M:%S")

    def _files(self, filter="*.txt"):
        files = []
        for root, _, filenames in os.walk(self.dir_path):
            for filename in fnmatch.filter(filenames, filter):
                files.append(join(root, filename))

        return sorted(files)

    def _repository(self, url):
        url_parse = urlparse(url.strip())
        return url_parse.path.replace(".git", "")

    def _blocks(self, content, re_start, re_stop=None, jump=0):
            blocks = []
            block = []
            init = False
            count = 0
            for line in content:
                count += 1
                if re.match(re_start, line) != None:
                    if block:
                        blocks.append(block)
                    init, count, block = True, 0, []

                elif re_stop and re.match(re_stop, line) != None:
                    if block:
                        blocks.append(block)
                    init, count, block = False, 0, []

                if init and count >= jump:
                    block.append(line)

            if init and block:
                blocks.append(block)

            return blocks

    def trim(self, value):
        return value.strip("\n\t\r ")

    def get_value_fields(self, data, fields, keys):
        ctx = {}
        for key in keys:
            field_pattern = fields[key]["pattern"]
            field_format = fields[key]["format"]
            value = self.get_value(data, field_pattern)
            if value and self.trim(value) != "":
                ctx[key] = field_format(self.trim(value))
            else:
                ctx[key] = ""

        return ctx

    def get_value(self, content, re_pattern):
        for line in content:
            resp = re.search(re_pattern, line)
            if resp:
                return resp.groups()[0]


class MhartNode(OfficialImages):

    def __init__(self, dir_path, language):
        log_debug(dir_path, language)

        self.hub = Hub(language)
        self.hub_join = Hub("node")
        self.language = language
        self.dir_path = dir_path
        self.files = self._files()
        self.data = self._load()
        self.repository = "/mhart/alpine-node"
        self.directory = ""
        self._parser()
        self._join()

    def _join(self):
        official_images = OfficialImages("tmp/manifest/node", "node")
        for image in official_images.hub.images:
            self.hub_join.add_image(image)
        for image in self.hub.images:
            self.hub_join.add_image(image)

    def _parser(self):
        for d in self.data:
            git_commit = d["filename"].split("_")[2]
            image = self.get_value(d["data"], r"^FROM\s+([-\w.:/]+)")
            version = self.get_value(d["data"], r"^ENV\s+VERSION=([\w.]+)")
            if not image or not version:
                continue

            tags = "{}-{}".format(version.strip("v"), image)
            image_kwargs = {
                "language": self.language,
                "git_datetime": d["git_datetime"],
                "tags": tags,
                "git_commit": git_commit,
                "directory": self.directory,
                "repository": self.repository,
                "filepath": d["filepath"],
            }
            image = Image(**image_kwargs)
            if image.is_valid():
                self.hub.add_image(image)

    def save(self):
        self.hub_join.save()

def log_debug(*args, **kwags):
    global DEBUG
    if DEBUG:
        if args and kwags:
            print(*args, kwags)
        elif args:
            print(*args)
        else:
            print(kwags)

def save_yml(file_name, data):
    yaml = YAML()
    yaml.indent(mapping=2, sequence=4, offset=2)
    with codecs.open(file_name, 'w', encoding='utf8') as f:
        yaml.dump(data, f)

    yaml_data = StringIO()
    with codecs.open(file_name, 'r+', encoding='utf8') as f:
        for line in f.readlines():
            if line.find("- version:") != -1 or line.find("- name:") != -1:
                yaml_data.write("\n")
            yaml_data.write(line)

        f.seek(0)
        f.write(yaml_data.getvalue())




if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usar no seguinte formato:\n python official-images.py <filepath> <language>\n")
        print("Exemplo:\n python official-images.py /tmp/all_versions_exported/python python True\n\n")
        exit(2)

    dir_path = sys.argv[1]
    language = sys.argv[2]
    DEBUG = False
    method_name = "save"

    if len(sys.argv) > 3:
        DEBUG = sys.argv[3] == 'True'

    if len(sys.argv) > 4:
        method_name = sys.argv[4]

    if language == "mhart":
        official_images = MhartNode(dir_path, language)
    elif language == "distributions":
        official_images = Distribution()
    else:
        official_images = OfficialImages(dir_path, language)

    if method_name != "":
        method = getattr(official_images, method_name)
        method()
