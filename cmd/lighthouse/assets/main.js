document.addEventListener('alpine:init', function() {
    // Reuse the same UUID across multiple calls to id() by appending a counter, yielding unique
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

    // Keep the storage key in sync with the inline script in views/layout/p/theme.go.html.
    Alpine.data('theme', function() {
        return {
            theme: localStorage.getItem('theme') ||
                (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'),

            themes: {
                light: { next: 'dark', icon: 'fa-sun-bright', label: 'Light theme' },
                dark: { next: 'light', icon: 'fa-moon', label: 'Dark theme' },
            },

            toggleTheme() {
                this.theme = this.themes[this.theme].next
                document.documentElement.dataset.theme = this.theme
                localStorage.setItem('theme', this.theme)
            },

            themeIcon() {
                return this.themes[this.theme].icon
            },

            themeLabel() {
                return this.themes[this.theme].label
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
        const SEARCH_DELAY = 200
        const PLACEHOLDER_CLASS = 'fa-solid fa-question'
        const STYLES = ['duotone', 'solid', 'brands']

        let controller = null
        let timer = null
        let isStale = false
        let pickWhenReady = false

        return {
            icon: null,
            isOpen: false,
            isLoading: false,
            results: {},
            // -1 means the search box is selected rather than a dropdown item.
            selectedIndex: -1,

            iconClass() {
                return this.hasIcon() ? this.iconToClass(...this.icon.split(':', 2)) : PLACEHOLDER_CLASS
            },

            hasIcon() {
                const [style, name] = (this.icon || '').split(':', 2)
                if (!STYLES.includes(style) || !name) {
                    return false
                }

                const probe = this.$refs.probe
                probe.className = `icon-probe ${this.iconToClass(style, name)}`
                return getComputedStyle(probe, '::before').content != 'none'
            },

            iconToClass(style, name) {
                return `fa-${style} fa-${name}`
            },

            searchTerm() {
                return (this.icon || '').split(':').pop()
            },

            async search() {
                this.cancelQueuedSearch()
                this.abortSearch()
                this.selectedIndex = -1

                if (!this.icon) {
                    isStale = false
                    pickWhenReady = false
                    this.results = {}
                    this.isOpen = false
                    return
                }

                const query = this.icon
                const current = new AbortController()
                controller = current
                this.isLoading = true

                try {
                    const url = '/api/icon/search?q=' + encodeURIComponent(query)
                    const response = await fetch(url, { signal: current.signal })
                    this.results = response.ok ? await response.json() : { failed: true }
                } catch (error) {
                    if (error.name == 'AbortError') {
                        return
                    }
                    this.results = { failed: true }
                } finally {
                    if (controller == current) {
                        controller = null
                        this.isLoading = false
                    }
                }

                isStale = this.icon != query
                this.isOpen = true

                if (pickWhenReady && !isStale) {
                    pickWhenReady = false
                    const [first] = this.options()
                    if (first) {
                        this.setIcon(first)
                    }
                }
            },

            queueSearch() {
                isStale = true
                this.cancelQueuedSearch()
                timer = setTimeout(() => {
                    timer = null
                    this.search()
                }, SEARCH_DELAY)
            },

            cancelQueuedSearch() {
                clearTimeout(timer)
                timer = null
            },

            abortSearch() {
                controller?.abort()
                controller = null
                this.isLoading = false
            },

            setIcon(icon) {
                this.icon = icon.style + ':' + icon.name
                this.closeDropdown()
                this.$refs.input.focus()
            },

            closeDropdown() {
                this.cancelQueuedSearch()
                isStale = false
                pickWhenReady = false
                this.abortSearch()
                this.isOpen = false
                this.selectedIndex = -1
            },

            options() {
                return (this.isOpen && this.results.icons) || []
            },

            moveSelectionUp() {
                this.selectedIndex = Math.max(this.selectedIndex - 1, -1)
            },

            moveSelectionDown() {
                if (!this.isOpen) {
                    this.search()
                    return
                }

                this.selectedIndex = Math.min(this.selectedIndex + 1, this.options().length - 1)
            },

            chooseSelection(event) {
                if (isStale && this.icon) {
                    event.preventDefault()
                    pickWhenReady = true
                    if (timer || !controller) {
                        this.search()
                    }
                    return
                }

                const options = this.options()
                if (options.length == 0) {
                    return
                }

                event.preventDefault()
                this.setIcon(options[Math.max(this.selectedIndex, 0)])
            },
        }
    })

    Alpine.data('filter', function() {
        return {
            filterSelection: null,

            filter() {
                const qs = new URLSearchParams(document.location.search)

                if (this.filterSelection) {
                    qs.set('f', this.filterSelection)
                } else {
                    qs.delete('f')
                }

                document.location.search = qs.toString()
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
