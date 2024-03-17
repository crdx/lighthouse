set quiet := true

PREFIX := 'cmd/lighthouse'

# create service e.g. fooer
service name:
    mkdir -pv {{ PREFIX }}/services/{{ name }}
    touch {{ PREFIX }}/services/{{ name }}/{{ name }}.go
    echo 'package {{ name }}' > {{ PREFIX }}/services/{{ name }}/{{ name }}.go

# create model e.g. foo
model name:
    mkdir -p db/repo/{{ name }}R
    touch db/repo/{{ name }}R/{{ name }}.go
    echo 'package {{ name }}R' > db/repo/{{ name }}R/{{ name }}.go
    touch queries/{{ name }}s.sql
    just make::migration 'create_{{ name }}_model'

# create controller e.g. foo
controller name:
    mkdir -pv {{ PREFIX }}/controllers/{{ name }}
    touch {{ PREFIX }}/controllers/{{ name }}/routes.go
    touch {{ PREFIX }}/controllers/{{ name }}/routes_test.go
    echo 'package {{ name }}' > {{ PREFIX }}/controllers/{{ name }}/routes.go
    echo 'package {{ name }}_test' > {{ PREFIX }}/controllers/{{ name }}/routes_test.go

# create migration e.g. create_foo_model
migration name:
    go run ./tools/make-migration "{{ name }}"
