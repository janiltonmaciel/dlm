class DockerfilerGen < Formula
    desc "Generator Dockerfile"
    homepage "https://github.com/janiltonmaciel/dockerfile-gen"
    url "https://github.com/janiltonmaciel/dockerfile-gen/archive/1.0.0.tar.gz"
    sha256 "67d9a16ebc0d7b979a5d95dd14ae7322f580b07bd3766bdf930bd7b75d0e53cc"
  
    depends_on "dep" => :build
    depends_on "go" => :build
  
    def install
      ENV["GOPATH"] = buildpath
  
      srcpath = buildpath/"src/github.com/janiltonmaciel/dockerfile-gen"
      srcpath.install buildpath.children
  
      cd srcpath do
        system "dep", "ensure", "-vendor-only"
        system "go", "build", "-ldflags", "-X main.buildVersion=#{version}", "-o", bin/"dockerfiler-gen", "main.go"
        prefix.install_metafiles
      end
    end
  
    test do
      assert_match version.to_s, shell_output("#{bin}/dockerfiler-gen --version")
    end
  end
  