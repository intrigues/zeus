function showToast() {
    var toastElList = [].slice.call(document.querySelectorAll('.toast'))
    var toastList = toastElList.map(function (toastEl) {
        // Creates an array of toasts (it only initializes them)
        return new bootstrap.Toast(toastEl) // No need for options; use the default options
    });
    toastList.forEach(toast => toast.show()); // This show them
}

function changeThemeToDark() {
    document.documentElement.setAttribute("data-theme", "dark")//set theme to light
    localStorage.setItem("data-theme", "dark") // save theme to local storage
}

const changeThemeToLight = () => {
    document.documentElement.setAttribute("data-theme", "light") // set theme light
    localStorage.setItem("data-theme", 'light') // save theme to local storage
}

function darkModeApply () {
    const toggleDarkMode = document.getElementById("toggleDarkMode");
    let theme = localStorage.getItem('data-theme'); // Retrieve saved them from local storage
    if (theme ==='dark'){
        changeThemeToDark()
        toggleDarkMode.classList.add("active")
    }else{
        changeThemeToLight()
        if (toggleDarkMode.classList.contains("active")) {
            toggleDarkMode.classList.remove("active")
        }
    } 
}
// Apply retrived them to the website
const toggleButton = document.getElementById("toggleDarkMode")
if (toggleButton) {
    toggleButton.addEventListener('click', () => {
        let toggleDarkMode = document.getElementById("toggleDarkMode");
        let theme = localStorage.getItem('data-theme'); // Retrieve saved them from local storage
        if (theme ==='dark'){
            changeThemeToLight()
            toggleDarkMode.checked = false
        }else{
            changeThemeToDark()
            toggleDarkMode.checked = true
        } 
    });
}


document.addEventListener('DOMContentLoaded', (event) => {
    // darkModeApply()
    let toggleDarkMode = document.getElementById("toggleDarkMode");
    let theme = localStorage.getItem('data-theme'); // Retrieve saved them from local storage
    if (theme ==='dark'){
        changeThemeToDark()
        if (toggleDarkMode) {
            toggleDarkMode.checked = true
        }
    }else{
        changeThemeToLight()
        if (toggleDarkMode) {
            toggleDarkMode.checked = false
        }
    } 
})