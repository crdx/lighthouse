document.addEventListener('DOMContentLoaded', function() {
    const burgers = document.querySelectorAll('.navbar-burger') || []

    burgers.forEach(el => {
        el.addEventListener('click', function() {
            const target = document.getElementById(el.dataset.target)
            el.classList.toggle('is-active')
            target.classList.toggle('is-active')
        })
    })
})
