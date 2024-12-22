// For the "Get Source Code" button
const get_code_buttons = document.querySelectorAll('.get-plus-btn');
get_code_buttons.forEach(button => {
    button.addEventListener('click', function() {
        const hiddenBox = this.nextElementSibling;
        hiddenBox.classList.toggle('show'); // Toggle 'show' class
    });
});

// For the "Terraform Plan" button
const terraform_plan_buttons = document.querySelectorAll('.tfplan-plus-btn');
terraform_plan_buttons.forEach(button => {
    button.addEventListener('click', function() {
        const hiddenBox = this.nextElementSibling;
        hiddenBox.classList.toggle('show'); // Toggle 'show' class
    });
});

// For the "Terraform Apply" button
const terraform_apply_buttons = document.querySelectorAll('.tfapply-plus-btn');
terraform_apply_buttons.forEach(button => {
    button.addEventListener('click', function() {
        const hiddenBox = this.nextElementSibling;
        hiddenBox.classList.toggle('show'); // Toggle 'show' class
    });
});
