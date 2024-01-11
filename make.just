# scaffold a new service (name is lowercase e.g. watcher)
@service name:
    mkdir -pv src/services/{{ name }}
    touch src/services/{{ name }}/{{ name }}.go
    echo 'package {{ name }}' > src/services/{{ name }}/{{ name }}.go

# scaffold a new model (name is lowercase e.g. device)
@model name:
    touch src/m/{{ name }}.go
    echo 'package m' > src/m/{{ name }}.go
    mkdir src/m/repo/{{ name }}R
    touch src/m/repo/{{ name }}R/{{ name }}.go
    echo 'package {{ name }}R' > src/m/repo/{{ name }}R/{{ name }}.go

# scaffold a new controller (name is lowercase e.g. device)
@controller name:
    mkdir -pv src/controllers/{{ name }}Controller
    touch src/controllers/{{ name }}Controller/routes.go
    touch src/controllers/{{ name }}Controller/routes_test.go
    echo 'package {{ name }}Controller' > src/controllers/{{ name }}Controller/routes.go
    echo 'package {{ name }}Controller_test' > src/controllers/{{ name }}Controller/routes_test.go

# scaffold a new migration (name is ProperCase with no punctuation)
[no-exit-message]
@migration name:
    go run src/cmd/mkmigration/main.go "{{ name }}"
