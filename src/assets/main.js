document.addEventListener('alpine:init', function() {
    // Save CPU cycles by reusing the same UUID across multiple calls to id(), yielding unique id
    // values that start with an alpha character.
    let prefix
    let counter = 1

    function id() {
        if (!prefix) {
            prefix = uuid()
        }
        return 'i_' + prefix + '_' + counter++
    }

    Alpine.data('id', function() {
        return {
            id: id()
        }
    })

    Alpine.data('nav', function() {
        return {
            navOpen: false,

            toggleNav() {
                this.navOpen = !this.navOpen
            },

            navClass() {
                return this.navOpen && 'is-active'
            },
        }
    })

    Alpine.data('dropdown', function() {
        return {
            dropdownOpen: false,

            closeDropdown() {
                this.dropdownOpen = false
            },

            toggleDropdown() {
                this.dropdownOpen = !this.dropdownOpen
            },

            dropdownClass() {
                return this.dropdownOpen && 'is-active'
            },
        }
    })

    Alpine.data('form', function() {
        return {
            submitForm() {
                this.$el.querySelector('form').requestSubmit()
            },
        }
    })

    Alpine.data('modal', function() {
        return {
            modalOpen: false,

            modalClass() {
                return this.modalOpen && 'is-active'
            },

            openModal() {
                this.modalOpen = true
            },

            closeModal() {
                this.modalOpen = false
            },
        }
    })

    Alpine.data('iconSearch', function() {
        return {
            icon: null,
            isOpen: false,
            results: [],
            // -1 means the search box is selected rather than a dropdown item.
            selectedIndex: -1,

            iconClass() {
                if (this.icon) {
                    let [style, name] = this.icon.split(':', 2)
                    if (style && name) {
                        return this.iconToClass(style, name)
                    }
                }
            },

            iconToClass(style, name) {
                return `fa-${style} fa-${name}`
            },

            async search() {
                this.selectedIndex = -1

                if (this.icon) {
                    this.results = (await (await fetch('/api/icon/search?q=' + this.icon)).json())
                }

                this.isOpen = !!this.icon
            },

            setIcon(icon) {
                this.icon = icon.style + ':' + icon.name
                this.closeDropdown()
            },

            closeDropdown() {
                this.isOpen = false
                this.$refs.input.focus()
            },

            moveSelectionUp() {
                if (this.selectedIndex == 0) {
                    this.selectedIndex = -1
                    this.$refs.input.focus()
                    return
                }

                if (this.selectedIndex > -1) {
                    this.selectedIndex--
                }

                this.dropdownItems()[this.selectedIndex].focus()
            },

            moveSelectionDown() {
                if (!this.isOpen) {
                    this.search()
                    return
                }

                if (this.selectedIndex < this.dropdownItems().length - 1) {
                    this.selectedIndex++
                }

                this.dropdownItems()[this.selectedIndex].focus()
            },

            dropdownItems() {
                return this.$refs.content.querySelectorAll('a')
            },
        }
    })
})

function uuid() {
    const a = new Uint8Array(16)
    crypto.getRandomValues(a)

    a[6] = (a[6] & 0x0f) | 0x40 // v4
    a[8] = (a[8] & 0x3f) | 0x80

    const segments = [
        a.subarray(0, 4),
        a.subarray(4, 6),
        a.subarray(6, 8),
        a.subarray(8, 10),
        a.subarray(10, 16),
    ]

    const f = bytes => [...bytes].map(b => b.toString(16).padStart(2, '0')).join('')
    return segments.map(f).join('')
}
