// jQuery-inspired aliases.
const $$ = s => document.querySelectorAll(s)
const $ = s => document.querySelector(s)

document.addEventListener('DOMContentLoaded', function() {
    initNav()
    initModals()
    initDropdowns()
    initForms()
})

function initNav() {
    $$('.navbar-burger').forEach(el => {
        el.addEventListener('click', function() {
            const target = $(el.dataset.target)
            el.classList.toggle('is-active')
            target.classList.toggle('is-active')
        })
    })
}

function initModals() {
    function open(el) {
        el.classList.add('is-active')
        $('html').classList.add('is-clipped')
    }

    function close(el) {
        el.classList.remove('is-active')
        $('html').classList.remove('is-clipped')
    }

    function closeAll() {
        $$('.modal').forEach(close)
    }

    $$('.js-modal-trigger').forEach(el => {
        const target = $(el.dataset.target)
        el.addEventListener('click', () => {
            open(target)
        })
    })

    $$('.modal-background, .modal-close, .modal-card-head .delete, .modal-card-foot .button').forEach(el => {
        const target = el.closest('.modal')
        el.addEventListener('click', () => {
            close(target)
        })
    })

    document.addEventListener('keydown', event => {
        if (event.code === 'Escape') {
            closeAll()
        }
    })
}

function initDropdowns() {
    $$('.js-dropdown-trigger').forEach(el => {
        const target = el.closest('.dropdown')
        el.addEventListener('click', () => {
            target.classList.toggle('is-active')
        })
    })
}

function initForms() {
    $$('.js-form-trigger').forEach(el => {
        el.addEventListener('click', () => {
            $(el.dataset.target).requestSubmit()
            return false
        })
    })
}

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

const id = (function() {
    // Save CPU cycles by reusing the same UUID across multiple calls to id(), yielding unique id
    // values that start with an alpha character.
    let prefix
    let counter = 1

    return function() {
        if (!prefix) {
            prefix = uuid()
        }
        return 'i_' + prefix + '_' + counter++
    }
})()
