"use strict";

const deleteBtns = document.querySelectorAll(".delete-post-btn");
const publishDraftBtns = document.querySelectorAll(".publish-draft-btn");

const deletePostEventHandler = (e) => {
  e.preventDefault();
  const postId = e.target.dataset.postId;
  const url = `/blog/dashboard/${postId}`;

  if (!confirm("Are you sure you want to delete the blog?")) return;

  fetch(url, {
    method: "DELETE",
  })
    .then((res) => res.json())
    .then((data) => {
      if (data["status"] == "success") {
        alert("successfully deleted post");
        window.location.href = window.location.href;
      } else {
        alert("failed to delete post");
      }
    });
};

const publishDraftPostEventHandler = (e) => {
  e.preventDefault();
  const postId = e.target.dataset.postId;
  const url = `/blog/dashboard/${postId}/toggle_publish`;

  fetch(url, {
    method: "PUT",
  })
    .then((res) => res.json())
    .then((data) => {
      if (data["status"] == "success") {
        let msg = "published";
        if (e.target.innerText === "[DRAFT]") {
          msg = "drafted";
        }
        alert(`successfully ${msg} the post`);
        window.location.href = window.location.href;
      } else {
        alert("failed to toggle publish");
      }
    });
};

deleteBtns.forEach((btn) =>
  btn.addEventListener("click", deletePostEventHandler)
);
publishDraftBtns.forEach((btn) =>
  btn.addEventListener("click", publishDraftPostEventHandler)
);
