const form = document.getElementById("addTaskForm");

form.addEventListener("submit", function (event) {
  event.preventDefault();

  const taskMatkul = document.getElementById("taskMatkul").value;
  let deadline = document.getElementById("deadline").value;

  if (deadline.length === 16) {
    deadline += ":00";
  }

  const taskData = {
    matkul: taskMatkul,
    deadline: deadline,
  };

  fetch("/add-task", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(taskData),
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.status === "success") {
        alert("Task added successfully");
        location.reload();
      }
    })
    .catch((error) => {
      console.error("Error:", error);
      alert("Failed to add task");
    });
});

function deleteTask(taskID) {
  fetch(`/delete-task?id=${taskID}`, {
    method: "DELETE",
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.status === "success") {
        alert("Task deleted successfully");
        location.reload();
      } else {
        alert("Failed to delete task");
      }
    })
    .catch((error) => {
      console.error("Error:", error);
      alert("Error deleting task");
    });
}

// Contoh membuka modal menggunakan Bootstrap 5
// Function to open the edit modal and populate the form with task data
function openEditModal(id, matkul, deadline) {
  // Set the task ID in the hidden input
  document.getElementById("editTaskID").value = id;

  // Populate the task data in the input fields
  document.getElementById("editTaskMatkul").value = matkul;
  document.getElementById("editDeadline").value = deadline;

  // Open the modal
  const modalElement = document.getElementById("editTaskModal"); // Fixed modal ID
  const modal = new bootstrap.Modal(modalElement);
  modal.show(); // Show the modal
}

// Function to update the task when the user clicks "Save changes"
function updateTask() {
  const taskID = document.getElementById("editTaskID").value;
  const taskMatkul = document.getElementById("editTaskMatkul").value;
  const deadline = document.getElementById("editDeadline").value;

  const taskData = {
    matkul: taskMatkul,
    deadline: deadline,
  };

  fetch(`/update-task?id=${taskID}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(taskData),
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.status === "success") {
        alert("Task updated successfully");
        location.reload();
      } else {
        alert("Failed to update task");
      }
    })
    .catch((error) => {
      console.error("Error:", error);
      alert("Error updating task");
    });
}

// Function to update the task status when checkbox is clicked
document.querySelectorAll(".statusCheckbox").forEach((checkbox) => {
  checkbox.addEventListener("change", function () {
    const taskID = this.getAttribute("data-id");
    // Kirim 1 jika checkbox dicentang, 0 jika tidak dicentang
    const completed = this.checked ? 1 : 0;

    fetch(`/update-status?id=${taskID}&completed=${completed}`, {
      method: "PUT",
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.status !== "success") {
          alert("Failed to update status");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
        alert("Failed to update status");
      });
  });
});
