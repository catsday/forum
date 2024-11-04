const titleInput = document.querySelector("#title");
const maxChars = 40;

if (titleInput) {
    titleInput.addEventListener("keydown", function(event) {
        if (this.value.length >= maxChars && !["Backspace", "Delete", "ArrowLeft", "ArrowRight"].includes(event.key)) {
            event.preventDefault();
            showAlert(`Title can have a maximum of ${maxChars} characters.`);
        }
    });

    titleInput.addEventListener("input", function() {
        if (this.value.length > maxChars) {
            this.value = this.value.slice(0, maxChars);
            showAlert(`Title can have a maximum of ${maxChars} characters.`);
        }
    });
}

function showAlert(message) {
    if (document.querySelector(".alert-box")) return;

    const alertBox = document.createElement("div");
    alertBox.classList.add("alert-box");
    alertBox.textContent = message;

    document.body.appendChild(alertBox);

    setTimeout(() => {
        alertBox.remove();
    }, 3000);
}

function filterByCategory(categoryID) {
    const buttons = document.querySelectorAll(".sidebar-item button");
    buttons.forEach(button => button.classList.remove("active"));

    const activeButton = document.querySelector(`button[onclick="filterByCategory(${categoryID})"]`);
    if (activeButton) activeButton.classList.add("active");

    window.location.href = `/?categoryID=${categoryID}`;
}

function resetFilter() {
    window.location.href = "/";
}

function filterMyPosts() {
    window.location.href = "/?myPosts=1";
}

function toggleVote(postID, voteType) {
    fetch("/toggle-vote", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `postID=${postID}&voteType=${voteType}`
    })
        .then(response => {
            if (response.ok) {
                window.location.reload();
            } else {
                alert("You must be logged in to vote.");
            }
        })
        .catch(error => console.error("Error:", error));
}

function filterLikedPosts() {
    window.location.href = "/?likedPosts=1";
}

const form = document.querySelector("form[action='/forum/create']");
if (form) {
    form.addEventListener("submit", function(event) {
        const categoryCheckboxes = document.querySelectorAll("input[name='categories']");
        const isCategorySelected = Array.from(categoryCheckboxes).some(checkbox => checkbox.checked);

        if (!isCategorySelected) {
            event.preventDefault();
            showAlert("Please select at least one category for your post.");
        }
    });
}
