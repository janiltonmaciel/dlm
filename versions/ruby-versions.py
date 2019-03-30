# -*- encoding: utf-8 -*-
import os
import yaml
import codecs
import requests
from bs4 import BeautifulSoup
from collections import OrderedDict


class RubyVersions():
    URL_RELEASES_YAML = "https://raw.githubusercontent.com/ruby/www.ruby-lang.org/master/_data/releases.yml"
    URL_RELEASE_NEWS_PAGE = "https://www.ruby-lang.org{}"
    URL_RELEASES_PAGE = "https://www.ruby-lang.org/en/downloads/releases/"

    def __init__(self):
        # self.file_path_html = "ruby-releases-page.html"
        self.file_path_releases_yaml = "ruby-releases.yml"
        self.file_path_versions_yaml = "ruby-versions.yml"
        self.rewrite = True

    def generate(self):
        self.download()
        self.load()
        self.parse()
        self.save()
            
    def download(self):
        if self.rewrite or not os.path.exists(self.file_path_releases_yaml):
            download_file(self.URL_RELEASES_YAML, self.file_path_releases_yaml)

    def load(self):
        with codecs.open(self.file_path_releases_yaml, 'r', encoding='utf8') as fp:
            self.ruby_releases = yaml.safe_load(fp)

    def parse(self):
        self.versions = []
        for release in self.ruby_releases:
            version_major = release["version"][0:3]
            version = release["version"]
            sha256 = release.get("sha256", {}).get("xz")
            if not sha256:
                sha256 = self.parse_release_page(version, self.URL_RELEASE_NEWS_PAGE.format(release["post"]))
                if not sha256:
                    continue
            
            do = UnsortableOrderedDict()
            do["version"] = str(version)
            do["version_major"] = float(version_major)
            do["sha256"] = str(sha256)
            self.versions.append(do)

    def parse_release_page(self, version, url):
        file_name = "{}.tar.xz".format(version)
        html = request_url(url)
        soup = BeautifulSoup(html, 'html.parser')
        data = soup.select("#content > ul > li")
        for li in data:
            if li.code and li.p and li.p.a and li.p.a['href'].endswith(file_name):
                return self._get_sha256(li.code.get_text())

    def _get_sha256(self, text):
        for d in text.split("\n"):
            if d.startswith("SHA256"):
                return d.replace("SHA256:", "").strip()

    def save(self):
        print("Save {}...".format(self.file_path_versions_yaml))
        if self.versions:
            with codecs.open(self.file_path_versions_yaml, 'w', encoding='utf8') as f:
                yaml.dump(self.versions, f, default_flow_style=False)

    # def download_html(self):
    #     if os.path.exists(self.file_path_html):
    #         with codecs.open(self.file_path_html, 'r', encoding='utf8') as fp:
    #             self.html = fp.read()
    #     else:
        
            # self.html = request_url(self.URL_RELEASES_PAGE)
    #     if self.rewrite:
    #         self.save(self.file_path_html, self.html)

    # def parse_releases(self):
    #     ruby_releases = []
    #     soup = BeautifulSoup(self.html, 'html.parser')
    #     data = soup.select(".release-list > tr")
    #     if data:
    #         for tr in data:
    #             tds = tr.find_all("td")
    #             if tds:
    #                 ruby_releases.append({
    #                     "version": tds[0].get_text().replace("Ruby", "").strip(),
    #                     "post": tds[2].a['href'].strip()
    #                 })
        
    #     self.ruby_releases = ruby_releases


class UnsortableList(list):
    def sort(self, *args, **kwargs):
        pass

class UnsortableOrderedDict(OrderedDict):
    def items(self, *args, **kwargs):
        return UnsortableList(OrderedDict.items(self, *args, **kwargs))

yaml.add_representer(UnsortableOrderedDict, yaml.representer.SafeRepresenter.represent_dict)


def download_file(url, file_path=None):
    local_filename = file_path or url.split('/')[-1]
    with requests.get(url, stream=True) as r:
        r.raise_for_status()
        with open(local_filename, 'wb') as f:
            for chunk in r.iter_content(chunk_size=512): 
                if chunk:
                    f.write(chunk)
    
    return local_filename

def request_url(url):
    with requests.get(url) as r:
        r.raise_for_status()
        return r.text
