let headerLinks = document.querySelectorAll('nav > a.nav-link')

let activeButton = "active"

headerLinks.forEach(function(item, index, arr) {

    item.addEventListener("click", function(e){


    if (!e.target.className.includes(activeButton)) {
        arr.forEach((el, ind) => {
            if (ind !== index) {
                el.classList.remove(activeButton)
            } else {
                el.classList.add(activeButton)
            }
        })
    }


    })
})