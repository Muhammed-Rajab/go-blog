"use strict";

const deleteBtns = document.querySelectorAll(".delete-post-btn");
const publishDraftBtns = document.querySelectorAll(".publish-draft-btn");

const deletePostEventHandler = (e) => {
  e.preventDefault();
  const postId = e.target.dataset.postId;
  const url = `/blog/dashboard/${postId}`;

  fetch(url, {
    method: "DELETE",
  })
    .then((res) => res.json())
    .then((data) => {
      if (data["status"] == "success") {
        alert("successfully delete post");
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
        alert("successfully toggled publish");
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
