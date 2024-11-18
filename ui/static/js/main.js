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

function toggleVote(postID, voteType) {
    fetch("/toggle-vote", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `postID=${postID}&voteType=${voteType}`
    })
        .then(response => {
            if (response.status === 401) {
                window.location.href = "/forum/login";
            } else if (response.ok) {
                window.location.reload();
            } else {
                alert("An error occurred while attempting to vote.");
            }
        })
        .catch(() => {
            alert("An error occurred while attempting to vote.");
        });
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

function filterByCategory(category) {
    switch (category) {
        case 1:
            window.location.href = "/forum/technology";
            break;
        case 2:
            window.location.href = "/forum/entertainment";
            break;
        case 3:
            window.location.href = "/forum/sports";
            break;
        case 4:
            window.location.href = "/forum/education";
            break;
        case 5:
            window.location.href = "/forum/health";
            break;
        default:
            window.location.href = "/";
    }
}

function filterLikedPosts() {
    window.location.href = "/forum/liked";
}

function filterMyPosts() {
    window.location.href = "/forum/posted";
}

function filterComments() {
    window.location.href = "/forum/commented";
}

function resetFilter() {
    window.location.href = "/";
}

function sortPosts(order) {
    let url = new URL(window.location.href);

    switch (order) {
        case 'newest':
            url.searchParams.set("sort", "desc");
            url.searchParams.delete("sortBy");
            break;
        case 'oldest':
            url.searchParams.set("sort", "asc");
            url.searchParams.delete("sortBy");
            break;
        case 'likes':
            url.searchParams.set("sortBy", "likes");
            url.searchParams.delete("sort");
            break;
        case 'comments':
            url.searchParams.set("sortBy", "comments");
            url.searchParams.delete("sort");
            break;
    }

    window.location.href = url.toString();
}




