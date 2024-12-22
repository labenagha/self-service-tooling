// For the "Last Successful Run" & "Last Failed Run" buttons
const getPlusButtons = document.querySelectorAll('.plus-btn, .second-plus-btn');

getPlusButtons.forEach(button => {
    button.addEventListener('click', function () {
        // Find the parent container of the button (either .box1 or .box2)
        const parentBox = this.closest('.box1, .box2');
        
        // Find the hidden box within this parent container
        const hiddenBox = parentBox.querySelector('.hidden-box1, .hidden-box2');
        
        if (hiddenBox) {
            hiddenBox.classList.toggle('show'); // Toggle the 'show' class
        } else {
            console.error('Hidden box not found for button:', this);
        }
    });
});


// For the "Get Source Code" button
// const get_code_buttons = document.querySelectorAll('.get-plus-btn');
// get_code_buttons.forEach(button => {
//     button.addEventListener('click', function() {
//         const hiddenBox = this.nextElementSibling;
//         hiddenBox.classList.toggle('show');
//     });
// });

// For the "Terraform Plan" button
const terraform_plan_buttons = document.querySelectorAll('.tfplan-plus-btn');
terraform_plan_buttons.forEach(button => {
    button.addEventListener('click', function() {
        const hiddenBox = this.nextElementSibling;
        hiddenBox.classList.toggle('show');
    });
});

// For the "Terraform Apply" button
const terraform_apply_buttons = document.querySelectorAll('.tfapply-plus-btn');
terraform_apply_buttons.forEach(button => {
    button.addEventListener('click', function() {
        const hiddenBox = this.nextElementSibling;
        hiddenBox.classList.toggle('show');
    });
});
