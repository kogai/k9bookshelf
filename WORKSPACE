workspace(
    name = "k9books",
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "207fad3e6689135c5d8713e5a17ba9d1290238f47b9ba545b63d9303406209c6",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.24.7/rules_go-v0.24.7.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.24.7/rules_go-v0.24.7.tar.gz",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "b85f48fa105c4403326e9525ad2b2cc437babaa6e15a3fc0b1dbab0ab064bc7c",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.2/bazel-gazelle-v0.22.2.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.2/bazel-gazelle-v0.22.2.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

go_repository(
    name = "com_github_99designs_gqlgen",
    importpath = "github.com/99designs/gqlgen",
    sum = "h1:haLTcUp3Vwp80xMVEg5KRNwzfUrgFdRmtBY8fuB8scA=",
    version = "v0.13.0",
)

go_repository(
    name = "com_github_agnivade_levenshtein",
    importpath = "github.com/agnivade/levenshtein",
    sum = "h1:n6qGwyHG61v3ABce1rPVZklEYRT8NFpCMrpZdBUbYGM=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_andreyvit_diff",
    importpath = "github.com/andreyvit/diff",
    sum = "h1:bvNMNQO63//z+xNgfBlViaCIJKLlCJ6/fmUseuG0wVQ=",
    version = "v0.0.0-20170406064948-c7f18ee00883",
)

go_repository(
    name = "com_github_arbovm_levenshtein",
    importpath = "github.com/arbovm/levenshtein",
    sum = "h1:jfIu9sQUG6Ig+0+Ap1h4unLjW6YQJpKZVmUzxsD4E/Q=",
    version = "v0.0.0-20160628152529-48b4e1c0c4d0",
)

go_repository(
    name = "com_github_burntsushi_toml",
    importpath = "github.com/BurntSushi/toml",
    sum = "h1:WXkYYl6Yr3qBf1K79EBnL4mak0OimBfB0XUf9Vl28OQ=",
    version = "v0.3.1",
)

go_repository(
    name = "com_github_cpuguy83_go_md2man_v2",
    importpath = "github.com/cpuguy83/go-md2man/v2",
    sum = "h1:EoUDS0afbrsXAZ9YQ9jdu/mZ2sXgT1/2yyNng4PGlyM=",
    version = "v2.0.0",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_dgryski_trifles",
    importpath = "github.com/dgryski/trifles",
    sum = "h1:fRzb/w+pyskVMQ+UbP35JkH8yB7MYb4q/qhBarqZE6g=",
    version = "v0.0.0-20200323201526-dd97f9abfb48",
)

go_repository(
    name = "com_github_go_chi_chi",
    importpath = "github.com/go-chi/chi",
    sum = "h1:uQNcQN3NsV1j4ANsPh42P4ew4t6rnRbJb8frvpp31qQ=",
    version = "v3.3.2+incompatible",
)

go_repository(
    name = "com_github_gogo_protobuf",
    importpath = "github.com/gogo/protobuf",
    sum = "h1:2jyBKDKU/8v3v2xVR2PtiWQviFUyiaGk2rpfyFT8rTM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_google_go_cmp",
    importpath = "github.com/google/go-cmp",
    sum = "h1:x95R7cp+rSeeqAMI2knLtQ0DKlaBhv2NrtrOvafPHRo=",
    version = "v0.5.3",
)

go_repository(
    name = "com_github_gorilla_context",
    importpath = "github.com/gorilla/context",
    sum = "h1:9oNbS1z4rVpbnkHBdPZU4jo9bSmrLpII768arSyMFgk=",
    version = "v0.0.0-20160226214623-1ea25387ff6f",
)

go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    sum = "h1:KOwqsTYZdeuMacU7CxjMNYEKeBvLbxW+psodrbcEa3A=",
    version = "v1.6.1",
)

go_repository(
    name = "com_github_gorilla_websocket",
    importpath = "github.com/gorilla/websocket",
    sum = "h1:+/TMaTYc4QFitKJxsQ7Yye35DkWvkdLcvGKqM+x0Ufc=",
    version = "v1.4.2",
)

go_repository(
    name = "com_github_hashicorp_golang_lru",
    importpath = "github.com/hashicorp/golang-lru",
    sum = "h1:YDjusn29QI/Das2iO9M0BHnIbxPeyuCHsjMW+lJfyTc=",
    version = "v0.5.4",
)

go_repository(
    name = "com_github_kr_pretty",
    importpath = "github.com/kr/pretty",
    sum = "h1:L/CwN0zerZDmRFUapSPitk6f+Q3+0za1rQkzVuMiMFI=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_kr_pty",
    importpath = "github.com/kr/pty",
    sum = "h1:VkoXIwSboBpnk99O/KFauAEILuNHv5DVFKZMBN/gUgw=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_kr_text",
    importpath = "github.com/kr/text",
    sum = "h1:45sCR5RtlFHMR4UwH9sdQ5TC8v0qDQCHnXt+kaKSTVE=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_logrusorgru_aurora",
    importpath = "github.com/logrusorgru/aurora",
    sum = "h1:bqDmpDG49ZRnB5PcgP0RXtQvnMSgIF14M7CBd2shtXs=",
    version = "v0.0.0-20200102142835-e9ef32dff381",
)

go_repository(
    name = "com_github_matryer_moq",
    importpath = "github.com/matryer/moq",
    sum = "h1:dkN4G4jtl0zpssB9OPCb8D+m4zTHdnZpi62HhlGnyTU=",
    version = "v0.1.4",
)

go_repository(
    name = "com_github_mattn_go_colorable",
    importpath = "github.com/mattn/go-colorable",
    sum = "h1:snbPLB8fVfU9iwbbo30TPtbLRzwWu6aJS6Xh4eaaviA=",
    version = "v0.1.4",
)

go_repository(
    name = "com_github_mattn_go_isatty",
    importpath = "github.com/mattn/go-isatty",
    sum = "h1:wuysRhFDzyxgEmMf5xjvJ2M9dZoWAXNNr5LSBS7uHXY=",
    version = "v0.0.12",
)

go_repository(
    name = "com_github_mitchellh_mapstructure",
    importpath = "github.com/mitchellh/mapstructure",
    sum = "h1:SzB1nHZ2Xi+17FP0zVQBHIZqvwRN9408fJO8h+eeNA8=",
    version = "v1.3.3",
)

go_repository(
    name = "com_github_opentracing_basictracer_go",
    importpath = "github.com/opentracing/basictracer-go",
    sum = "h1:YyUAhaEfjoWXclZVJ9sGoNct7j4TVk7lZWlQw5UXuoo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_opentracing_opentracing_go",
    importpath = "github.com/opentracing/opentracing-go",
    sum = "h1:3jA2P6O1F9UOrWVpwrIo17pu01KWvNWg4X946/Y5Zwg=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_pkg_errors",
    importpath = "github.com/pkg/errors",
    sum = "h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
    version = "v0.9.1",
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    importpath = "github.com/pmezard/go-difflib",
    sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_rs_cors",
    importpath = "github.com/rs/cors",
    sum = "h1:G9tHG9lebljV9mfp9SNPDL36nCDxmo3zTlAf1YgvzmI=",
    version = "v1.6.0",
)

go_repository(
    name = "com_github_russross_blackfriday_v2",
    importpath = "github.com/russross/blackfriday/v2",
    sum = "h1:JIOH55/0cWyOuilr9/qlrm0BSXldqnqwMsf35Ld67mk=",
    version = "v2.1.0",
)

go_repository(
    name = "com_github_sergi_go_diff",
    importpath = "github.com/sergi/go-diff",
    sum = "h1:we8PVUC3FE2uYfodKH/nBHMSetSfHDR6scGdBi+erh0=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_shurcool_httpfs",
    importpath = "github.com/shurcooL/httpfs",
    sum = "h1:SWV2fHctRpRrp49VXJ6UZja7gU9QLHwRpIPBN89SKEo=",
    version = "v0.0.0-20171119174359-809beceb2371",
)

go_repository(
    name = "com_github_shurcool_sanitized_anchor_name",
    importpath = "github.com/shurcooL/sanitized_anchor_name",
    sum = "h1:PdmoCO6wvbs+7yrJyMORt4/BmY5IYyJwS/kOiWx8mHo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_shurcool_vfsgen",
    importpath = "github.com/shurcooL/vfsgen",
    sum = "h1:JJV9CsgM9EC9w2iVkwuz+sMx8yRFe89PJRUrv6hPCIA=",
    version = "v0.0.0-20180121065927-ffb13db8def0",
)

go_repository(
    name = "com_github_stretchr_objx",
    importpath = "github.com/stretchr/objx",
    sum = "h1:4G4v2dO3VZwixGIRoQ5Lfboy6nUhCyYzaqnIAPPhYs4=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    sum = "h1:2E4SXV/wtOkTonXsotYi4li6zVWxYlZuYNCXe9XRJyk=",
    version = "v1.4.0",
)

go_repository(
    name = "com_github_urfave_cli_v2",
    importpath = "github.com/urfave/cli/v2",
    sum = "h1:qph92Y649prgesehzOrQjdWyxFOp/QVM+6imKHad91M=",
    version = "v2.3.0",
)

go_repository(
    name = "com_github_vektah_dataloaden",
    importpath = "github.com/vektah/dataloaden",
    sum = "h1:ZfVN2QD6swgvp+tDqdH/OIT/wu3Dhu0cus0k5gIZS84=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_vektah_gqlparser_v2",
    importpath = "github.com/vektah/gqlparser/v2",
    sum = "h1:uiKJ+T5HMGGQM2kRKQ8Pxw8+Zq9qhhZhz/lieYvCMns=",
    version = "v2.1.0",
)

go_repository(
    name = "com_github_yamashou_gqlgenc",
    importpath = "github.com/Yamashou/gqlgenc",
    sum = "h1:VDHpRaUiijHN8s1KbD8ntf1MBofbp84gCEH+cYrQ/uA=",
    version = "v0.0.0-20201118110422-c5e4e29019a7",
)

go_repository(
    name = "com_github_yuin_goldmark",
    importpath = "github.com/yuin/goldmark",
    sum = "h1:ruQGxdhGHe7FWOJPT0mKs5+pD2Xs1Bm/kdGlHO04FmM=",
    version = "v1.2.1",
)

go_repository(
    name = "com_sourcegraph_sourcegraph_appdash",
    importpath = "sourcegraph.com/sourcegraph/appdash",
    sum = "h1:d2maSb13hr/ArmfK3rW+wNUKKfytCol7W1/vDHxMPiE=",
    version = "v0.0.0-20180110180208-2cc67fd64755",
)

go_repository(
    name = "com_sourcegraph_sourcegraph_appdash_data",
    importpath = "sourcegraph.com/sourcegraph/appdash-data",
    sum = "h1:e1sMhtVq9AfcEy8AXNb8eSg6gbzfdpYhoNqnPJa+GzI=",
    version = "v0.0.0-20151005221446-73f23eafcf67",
)

go_repository(
    name = "in_gopkg_check_v1",
    importpath = "gopkg.in/check.v1",
    sum = "h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=",
    version = "v1.0.0-20190902080502-41f04d3bba15",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    sum = "h1:clyUAQHOM3G0M3f5vQj7LuJrETvjVot3Z5el9nffUtU=",
    version = "v2.3.0",
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    sum = "h1:psW17arqaxU48Z5kZ0CQnkZWQJsqcURM6tKiBApRjXI=",
    version = "v0.0.0-20200622213623-75b288015ac9",
)

go_repository(
    name = "org_golang_x_mod",
    importpath = "golang.org/x/mod",
    sum = "h1:RM4zey1++hCTbCVQfnWeKs9/IEsaBLA8vTkd0WVtmH4=",
    version = "v0.3.0",
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    sum = "h1:IX6qOQeG5uLjB/hjjwjedwfjND0hgjPMMyO1RoIXQNI=",
    version = "v0.0.0-20201021035429-f5854403a974",
)

go_repository(
    name = "org_golang_x_sync",
    importpath = "golang.org/x/sync",
    sum = "h1:SQFwaSi55rU7vdNs9Yr0Z324VNlrF+0wMqRXT4St8ck=",
    version = "v0.0.0-20201020160332-67f06af15bc9",
)

go_repository(
    name = "org_golang_x_sys",
    importpath = "golang.org/x/sys",
    sum = "h1:+Nyd8tzPX9R7BWHguqsrbFdRx3WQ/1ib8I44HXV5yTA=",
    version = "v0.0.0-20200930185726-fdedc70b468f",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    sum = "h1:cokOdA+Jmi5PJGXLlLllQSgYigAEfHXJAERHVMaCc2k=",
    version = "v0.3.3",
)

go_repository(
    name = "org_golang_x_tools",
    importpath = "golang.org/x/tools",
    sum = "h1:z5+u0pkAUPUWd3taoTialQ2JAMo4Wo1Z3L25U4ZV9r0=",
    version = "v0.0.0-20201121010211-780cb80bd7fb",
)

go_repository(
    name = "org_golang_x_xerrors",
    importpath = "golang.org/x/xerrors",
    sum = "h1:go1bK/D/BFZV2I8cIQd1NKEZ+0owSTG1fDTci4IqFcE=",
    version = "v0.0.0-20200804184101-5ec99f83aff1",
)

go_repository(
    name = "com_github_gomarkdown_markdown",
    importpath = "github.com/gomarkdown/markdown",
    sum = "h1:/Y3B7hM9H3TOWPhe8eWGBGS4r09pjvS5Z0uoPADyjmU=",
    version = "v0.0.0-20201113031856-722100d81a8e",
)

go_repository(
    name = "com_github_mattn_go_runewidth",
    importpath = "github.com/mattn/go-runewidth",
    sum = "h1:3tS41NlGYSmhhe/8fhGRzc+z3AYCw1Fe1WAyLuujKs0=",
    version = "v0.0.8",
)

go_repository(
    name = "com_github_mattn_godown",
    importpath = "github.com/mattn/godown",
    sum = "h1:UeXqxG2wTENf2BataGYSoX3jUrJ5PI0ceUzdL1L+BPI=",
    version = "v0.0.0-20201027140031-2c7783b24de7",
)

go_repository(
    name = "org_golang_dl",
    importpath = "golang.org/dl",
    sum = "h1:jeP6FgaSLNTMP+Yri3qjlACywQLye+huGLmNGhBzm6k=",
    version = "v0.0.0-20190829154251-82a15e2f2ead",
)

go_rules_dependencies()

go_register_toolchains()

gazelle_dependencies()

# WORKSPACE can't use `select` embedded function
# https://docs.bazel.build/versions/master/configurable-attributes.html#why-doesnt-select-work-with-bind
http_file(
    name = "theme_darwin",
    executable = True,
    sha256 = "4b734d08db62e001adcf8ac5607414aa9cb6cfb410faa8f53b02f9f5c8a5733c",
    urls = ["https://shopify-themekit.s3.amazonaws.com/v1.0.2/darwin-amd64/theme"],
)

http_file(
    name = "theme_linux",
    executable = True,
    sha256 = "21a058bc47555ec343088c9322187bb0ad57e4055454f509c0751c092af6b12b",
    urls = ["https://shopify-themekit.s3.amazonaws.com/v1.0.2/linux-amd64/theme"],
)
