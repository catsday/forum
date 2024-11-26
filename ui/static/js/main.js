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

function openModal(modalId) {
    document.getElementById(modalId).style.display = 'block';
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
}

function openUserProfile(button) {
    document.getElementById('profile-username').textContent = button.getAttribute('data-user-username');
    document.getElementById('profile-email').textContent = button.getAttribute('data-user-email');
    document.getElementById('profile-id').textContent = button.getAttribute('data-user-id');
    document.getElementById('profile-postcount').textContent = button.getAttribute('data-user-postcount');
    document.getElementById('profile-commentcount').textContent = button.getAttribute('data-user-commentcount');
    document.getElementById('profile-likedposts').textContent = button.getAttribute('data-user-likedposts');
    document.getElementById('profile-dislikedposts').textContent = button.getAttribute('data-user-dislikedposts');
    document.getElementById('profile-ratioposts').textContent = button.getAttribute('data-user-ratioposts');

    openModal('view-profile-modal');
}


function toggleBan(userID, button) {
    const isBanned = button.dataset.isBanned === "true";
    const action = isBanned ? "unban" : "ban";

    if (confirm(`Are you sure you want to ${action} this user?`)) {
        fetch(`/forum/toggle-ban?userID=${userID}`, { method: 'POST' })
            .then(response => {
                if (response.ok) {
                    alert(`User has been ${action}ned successfully.`);
                    button.textContent = isBanned ? "Ban" : "Unban";
                    button.dataset.isBanned = !isBanned;
                } else {
                    alert("Failed to update user status.");
                }
            });
    }
}

document.addEventListener("DOMContentLoaded", () => {
    const banButtons = document.querySelectorAll(".ban-button");

    banButtons.forEach((button) => {
        button.addEventListener("click", () => {
            const userID = button.dataset.userId;
            const action = button.textContent.trim();
            const confirmMessage = action === "Ban"
                ? "Are you sure you want to ban this user?"
                : "Are you sure you want to unban this user?";

            showAlertWithConfirmationBan(confirmMessage, () => {
                fetch(`/forum/toggle-ban?userID=${userID}`, { method: "POST" })
                    .then((response) => {
                        if (response.ok) {
                            button.textContent = action === "Ban" ? "Unban" : "Ban";
                            const resultMessage = action === "Ban"
                                ? "User has been banned."
                                : "User has been unbanned.";
                            showAlertBan(resultMessage);
                        } else {
                            showAlertBan("Failed to update ban status. Please try again.");
                        }
                    })
                    .catch((error) => {
                        console.error("Error updating ban status:", error);
                        showAlertBan("An error occurred. Please try again.");
                    });
            });
        });
    });
});

function showAlertBan(message) {
    if (document.querySelector(".alert-box-ban")) return;

    const overlay = document.createElement("div");
    overlay.classList.add("alert-overlay-ban");

    const alertBox = document.createElement("div");
    alertBox.classList.add("alert-box-ban");
    alertBox.textContent = message;

    document.body.appendChild(overlay);
    document.body.appendChild(alertBox);

    setTimeout(() => {
        alertBox.remove();
        overlay.remove();
    }, 1000);
}

function showAlertWithConfirmationBan(message, onConfirm) {
    document.body.style.overflow = "hidden";

    const overlay = document.createElement("div");
    overlay.classList.add("alert-overlay-ban");

    const alertBox = document.createElement("div");
    alertBox.classList.add("alert-box-ban", "alert-confirm-ban");

    const alertMessage = document.createElement("p");
    alertMessage.textContent = message;

    const buttonGroup = document.createElement("div");
    buttonGroup.classList.add("alert-button-group-ban");

    const confirmButton = document.createElement("button");
    confirmButton.textContent = "Yes";
    confirmButton.classList.add("alert-confirm-button-ban");
    confirmButton.addEventListener("click", () => {
        document.body.style.overflow = "";
        overlay.remove();
        alertBox.remove();
        onConfirm();
    });

    const cancelButton = document.createElement("button");
    cancelButton.textContent = "No";
    cancelButton.classList.add("alert-cancel-button-ban");
    cancelButton.addEventListener("click", () => {
        document.body.style.overflow = "";
        overlay.remove();
        alertBox.remove();
    });

    buttonGroup.appendChild(confirmButton);
    buttonGroup.appendChild(cancelButton);
    alertBox.appendChild(alertMessage);
    alertBox.appendChild(buttonGroup);

    document.body.appendChild(overlay);
    document.body.appendChild(alertBox);
}

function toggleCommentVote(commentID, voteType) {
    fetch("/toggle-comment-vote", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `commentID=${commentID}&voteType=${voteType}`
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
