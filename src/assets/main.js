// jQuery-inspired aliases.
const $$ = s => document.querySelectorAll(s)
const $ = s => document.querySelector(s)

document.addEventListener('DOMContentLoaded', function() {
    initNav()
    initModals()
    initDropdowns()
    initForms()
    initTabs()
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

function initTabs() {
    $$('.js-tab-trigger').forEach(el => {
        const target = $(el.dataset.target)
        el.addEventListener('click', () => {
            $$('.js-tab-trigger').forEach(el => el.classList.remove('is-active'))
            el.classList.add('is-active')

            $$('.js-tab').forEach(el => el.style.display = 'none')
            target.style.display = 'block'
        })
    })
}
