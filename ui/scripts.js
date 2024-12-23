// // For the "Last Successful Run" & "Last Failed Run" buttons
// const getPlusButtons = document.querySelectorAll('.plus-btn, .second-plus-btn');

// getPlusButtons.forEach(button => {
//     button.addEventListener('click', function () {
//         // Find the parent container of the button (either .box1 or .box2)
//         const parentBox = this.closest('.box1, .box2');
        
//         // Find the hidden box within this parent container
//         const hiddenBox = parentBox.querySelector('.hidden-box1, .hidden-box2');
        
//         if (hiddenBox) {
//             hiddenBox.classList.toggle('show'); // Toggle the 'show' class
//         } else {
//             console.error('Hidden box not found for button:', this);
//         }
//     });
// });

// // "Get Source Code" button redirects to login
// // document.querySelector('.get-plus-btn').addEventListener('click', () => {
// //     window.location.href = '/login'; // Redirect to login endpoint
// // });

// // Terraform Plan functionality
// document.querySelector('.tfplan-plus-btn').addEventListener('click', async () => {
//     const outputBox = document.querySelector('.box4 .hidden-box');
//     outputBox.innerHTML = "Running Terraform Plan...";
//     outputBox.classList.add('show');

//     try {
//         const response = await fetch('/api/terraform/plan'); // Backend endpoint
//         const result = await response.text();
//         outputBox.innerHTML = `<pre>${result}</pre>`;
//     } catch (error) {
//         outputBox.innerHTML = `Error: ${error.message}`;
//     }
// });

// // Terraform Apply functionality
// document.querySelector('.tfapply-plus-btn').addEventListener('click', async () => {
//     const outputBox = document.querySelector('.box5 .hidden-box');
//     outputBox.innerHTML = "Running Terraform Apply...";
//     outputBox.classList.add('show');

//     try {
//         const response = await fetch('/api/terraform/apply'); // Backend endpoint
//         const result = await response.text();
//         outputBox.innerHTML = `<pre>${result}</pre>`;
//     } catch (error) {
//         outputBox.innerHTML = `Error: ${error.message}`;
//     }
// });