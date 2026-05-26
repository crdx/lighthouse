# list release tags
tags:
    git tag -l --sort=-creatordate

# extract release notes for a tag from CHANGELOG.md
notes tag:
    #!/usr/bin/env python3
    import re
    version = '{{ tag }}'.lstrip('v')
    lines, found = [], False
    for line in open('CHANGELOG.md'):
        if re.match(r'^## \[%s\]' % re.escape(version), line):
            found = True
            continue
        if found and line.startswith('## '):
            break
        if found:
            lines.append(line.rstrip('\n'))
    print('\n'.join(lines).strip())

# push commits and a specific tag
push tag:
    git push
    git push origin "{{ tag }}"
