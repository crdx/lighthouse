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
