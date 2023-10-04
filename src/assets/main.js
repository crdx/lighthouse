document.addEventListener('DOMContentLoaded', function() {
    const burgers = document.querySelectorAll('.navbar-burger') || []

    burgers.forEach(function(el) {
        el.addEventListener('click', function() {
            const target = document.getElementById(el.dataset.target)
            el.classList.toggle('is-active')
            target.classList.toggle('is-active')
        })
    });

    const deleteButtons = document.querySelectorAll('.notification .delete') || []

    deleteButtons.forEach(function(deleteNode) {
        deleteNode.addEventListener('click', function() {
            const notification = deleteNode.parentNode
            notification.parentNode.removeChild(notification)
        })
    })
})
