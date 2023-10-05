function findParentByClass(node, className) {
  let current = node

  while (current !== null) {
    if (current.classList && current.classList.contains(className)) {
      return current
    }
    current = current.parentNode
  }

  return null
}

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
            const notification = findParentByClass(deleteNode, 'container')
            notification.parentNode.removeChild(notification)
        })
    })
})
