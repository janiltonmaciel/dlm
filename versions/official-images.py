# -*- encoding: utf-8 -*-
import os
import re
import fnmatch
import sys
from datetime import datetime
# from distutils.version import StrictVersion, LooseVersion
from natsort import natsorted

from os.path import dirname, realpath, join, exists, basename
import codecs
import requests
from collections import OrderedDict
from unicodedata import normalize
import unicodedata

try:
    from urllib.parse import urlparse
except ImportError:
    from urlparse import urlparse


class Hub():

    def __init__(self):
        self.images = []

    def add_image(self, image):
        # print(image.version, image.plataform, image.release, image.stags)

        for data in self.images:
            inter = image.stags.intersection(data.stags)
            if len(inter) > 0:
                if image.git_datetime < data.git_datetime:
                    for tag in inter:
                        image.stags.remove(tag)
                else:
                    for tag in inter:
                        data.stags.remove(tag)

        if len(image.stags) > 0 :
            # print(image.version, image.is_prerelease)
            # log_debug(image.version, image.plataform, image.release, image.stags)
            # print(image.git_datetime, image.version)
            # print(image.version, image.plataform, image.release, image.stags, image.url_repository)
            self.images.append(image)

    def versions(self):
        images_versions = {image.version:image for image in self.images}
        versions = list(images_versions.keys())

        versions = natsorted(versions)
        for v in versions:
            print(v)


# prerelease
class Image():
    PLATAFORMS = {
        # DEBIAN
        # "buildpack-deps": {"plataform": "Debian", "release": "Debian XXXX"},

        "buildpack-deps:stretch": {"plataform": "Debian", "release": "Debian 9 (stretch)"},
        "buildpack-deps:stretch-scm": {"plataform": "Debian", "release": "Debian 9 (stretch)"},
        "buildpack-deps:jessie": {"plataform": "Debian", "release": "Debian 8 (jessie)"},
        "buildpack-deps:jessie-scm": {"plataform": "Debian", "release": "Debian 8 (jessie)"},
        
        "debian:stretch": {"plataform": "Debian", "release": "Debian 9 (slim-stretch)"},
        "debian:jessie": {"plataform": "Debian", "release": "Debian 8 (slim-jessie)"},
        "debian:wheezy": {"plataform": "Debian", "release": "Debian 7 (slim-wheezy)"},
        "debian:jessie-slim": {"plataform": "Debian", "release": "Debian 8 (jessie-slim)"},
        "debian:stretch-slim": {"plataform": "Debian", "release": "Debian 8 (stretch-slim)"},
        
        "slim-stretch": {"plataform": "Debian", "release": "Debian 9 (slim-stretch)"},
        "slim-jessie": {"plataform": "Debian", "release": "Debian 8 (slim-jessie)"},

        "jessie-slim": {"plataform": "Debian", "release": "Debian 8 (jessie-slim)"},
        "stretch-slim": {"plataform": "Debian", "release": "Debian 8 (stretch-slim)"},
        
        "stretch": {"plataform": "Debian", "release": "Debian 9 (stretch)"},
        "jessie": {"plataform": "Debian", "release": "Debian 8 (jessie)"},
        "wheezy": {"plataform": "Debian", "release": "Debian 7 (wheezy)"},

        # ALPINE
        "alpine": {"plataform": "Alpine", "release": None},
    }

    RE_PRERELEASE = re.compile(r"(rc|beta|a|b)+\d+$")
    RE_IGNORE = re.compile(r"(-preview|onbuild|windowsservercore|nanoserver|-cross|chakracore)")
    RE_FROM_TAGS = re.compile(r"^FROM\s+([-\w.:/]+)")
    RE_RELEASE = re.compile(r"([\.\d]+)$")
    RE_VERSIONS = {
        "ruby": re.compile(r"^(([.\d]+)(-p\d+)?)"),
        "python": re.compile(r"^(([.\d]+)(\w\d+)?)"),
        "golang": re.compile(r"^(([.\d]+)((beta|rc)\d+)?)"),
        "default": re.compile(r"^([.\d]+)"),
    }

    def __init__(self, language, git_datetime, tags, directory, git_commit, repository):
        self.language = language
        self.tags = tags
        self.stags = set(map(str.strip, tags.split(",")))
        self.git_datetime = git_datetime
        self.git_commit = git_commit
        self.directory = directory
        self.repository = repository
        self.re_versions = self.RE_VERSIONS.get(self.language) or self.RE_VERSIONS["default"]
        self.url_repository = self.format_url()

        self.plataform = None
        self.release = None
        self.version = None
        self.major_version = None
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
            log_debug("INVALID", self.tags)
            self.valid = False
            return

        self.plataform, self.release = self.find_plataform_release(self.tags)

        if not self.plataform or not self.release:
            self.plataform, self.release = self.get_plataform_release()
            if not self.plataform or not self.release:
                log_debug("INVALID", self.tags, self.url_repository)
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
        return version, version[0:3]

    def find_plataform_release(self, tags):
        plataform = None
        release = None
        stags = set(map(str.strip, tags.split(",")))
        for key, value in self.PLATAFORMS.items():
            if tags.find(key) != -1:
                plataform = value["plataform"]
                release = value["release"] or self._release(stags)
                break

        if not plataform:
            log_debug("Plataforma não encontrada ", tags)

        return plataform, release

    def get_plataform_release(self):
        text = self.download(self.url_repository)
        resp = self.RE_FROM_TAGS.search(text)
        if resp:
            tags = resp.groups()[0]
            log_debug("get_plataform_release:", tags)
            return self.find_plataform_release(tags)

        return None, None

    def _release(self, stags):
        release = None
        for tag in stags:
            resp = self.RE_RELEASE.search(tag)
            if resp :
                release = resp.groups()[0]
                break

        return release

    def is_valid(self):
        return not self.ignore and self.valid

    @property
    def is_prerelease(self):
        return self.RE_PRERELEASE.search(self.version) != None

    def request_url(self, url):
        with requests.get(url) as r:
            r.raise_for_status()
            return r.text

    def download(self, url):
        tmp_filepath = "/tmp/dockerfile-gen/{}/{}".format(self.language, self.slugify(url))
        if not os.path.exists(tmp_filepath):
            self.download_file(url, tmp_filepath)

        return self.load(tmp_filepath)

    def download_file(self, url, file_path):
        dir_path = dirname(file_path)
        if not os.path.exists(dir_path):
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
        return "{raw_git_url}{repository}/{commit}/{directory}/Dockerfile".format(
            raw_git_url=raw_git_url, 
            repository=self.repository, 
            commit=self.git_commit, 
            directory=self.directory)

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

        self.hub = Hub()
        self.language = language
        self.dir_path = dir_path
        self.files = self._files()
        self.data = self._load()
        self._parser()

    def print_versions(self):
        return self.hub.versions()

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
                    "data": fp.readlines()
                })

        return data

    def _git_datetime(self, filepath):
        filename = basename(filepath)
        dt = filename.split("_")[1]
        return datetime.strptime(dt, "%Y-%m-%dT%H:%M:%S")

    def _files(self):
        files = []
        for root, _, filenames in os.walk(self.dir_path):
            for filename in fnmatch.filter(filenames, '*.txt'):
                files.append(os.path.join(root, filename))

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


def log_debug(*args, **kwags):
    global DEBUG
    if DEBUG:
        if args and kwags:
            print("DEBUG: ", args, kwags)
        elif args:
            print("DEBUG: ", args)
        else:
            print("DEBUG: ", kwags)

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usar no seguinte formato:\n python official-images.py <filepath> <language>\n")
        print("Exemplo:\n python official-images.py /tmp/all_versions_exported/python python True\n\n")
        exit(2)

    dir_path = sys.argv[1]
    language = sys.argv[2]
    DEBUG = sys.argv[3] == 'True'
    OfficialImages(dir_path, language).print_versions()