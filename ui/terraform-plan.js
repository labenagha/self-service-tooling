// terraform-plan.js

// Wait until the DOM is fully loaded
document.addEventListener("DOMContentLoaded", function () {
    // Grab references to the container and the buttons
    const variablesContainer = document.getElementById("variables-container");
    const addVariableBtn = document.getElementById("add-variable-button");
    const saveButton = document.getElementById("save-button");
    const savedVariablesContainer = document.getElementById("saved-variables-container");

    // Function to create a new variable block
    function createVariableBlock(variable = {}) {
        const variableBlock = document.createElement("div");
        variableBlock.classList.add("variable-block");

        variableBlock.innerHTML = `
            <button class="remove-variable-btn">Remove</button>
            <div class="input-row">
                <div class="key-value-group">
                    <label class="label" for="key">Key</label>
                    <input type="text" class="input-text" name="key" placeholder="key" value="${variable.key || ''}" required />

                    <label class="label" for="value">Value</label>
                    <input type="text" class="input-text" name="value" placeholder="value" value="${variable.value || ''}" required />
                </div>

                <div class="checkbox-group">
                    <label class="checkbox-label">
                        <input type="checkbox" name="hcl" ${variable.hcl ? 'checked' : ''} />
                        HCL
                    </label>
                    <label class="checkbox-label">
                        <input type="checkbox" name="sensitive" ${variable.sensitive ? 'checked' : ''} />
                        Sensitive
                    </label>
                </div>
            </div>

            <div class="description-row">
                <label class="label" for="description">Description (Optional)</label>
                <textarea class="textarea" name="description" placeholder="description (optional)">${variable.description || ''}</textarea>
            </div>
        `;

        // Add event listener to the remove button
        const removeBtn = variableBlock.querySelector('.remove-variable-btn');
        removeBtn.addEventListener('click', function () {
            variablesContainer.removeChild(variableBlock);
            // Hide Save button if no variables remain
            if (variablesContainer.children.length === 0) {
                saveButton.classList.add('hidden');
            }
        });

        return variableBlock;
    }

    // Function to load saved variables from localStorage (if any)
    function loadSavedVariables() {
        const savedVariables = localStorage.getItem("terraformVariables");
        if (savedVariables) {
            try {
                const variables = JSON.parse(savedVariables);
                variables.forEach((variable) => {
                    const variableBlock = createVariableBlock(variable);
                    variablesContainer.appendChild(variableBlock);
                });

                // Show the Save button if variables are loaded
                if (variables.length > 0) {
                    saveButton.classList.remove("hidden");
                }

                // Display saved variables below
                displaySavedVariables(variables);
            } catch (error) {
                console.error("Error parsing saved variables:", error);
            }
        } else {
            savedVariablesContainer.innerHTML = "<p>No variables saved.</p>";
        }
    }

    function displaySavedVariables(variables) {
        // Clear existing display
        savedVariablesContainer.innerHTML = '';
    
        if (variables.length === 0) {
            savedVariablesContainer.innerHTML = "<p>No variables saved.</p>";
            return;
        }
    
        // Create the table structure
        const table = document.createElement('table');
        table.classList.add('variables-table');
    
        // Create the table header
        const headerRow = `
            <thead>
                <tr>
                    <th>Key</th>
                    <th>Value</th>
                    <th>Category</th>
                    <th>Actions</th>
                </tr>
            </thead>
        `;
        table.innerHTML = headerRow;
    
        // Create the table body
        const tbody = document.createElement('tbody');
        variables.forEach((variable, index) => {
            const category = variable.hcl ? 'HCL' : variable.sensitive ? 'Sensitive' : 'env';
            const row = `
                <tr>
                    <td>${escapeHtml(variable.key)}</td>
                    <td>${variable.sensitive ? 'Sensitive - write only' : escapeHtml(variable.value)}</td>
                    <td>${category}</td>
                    <td>
                        <div class="actions">
                            <button class="edit-btn" data-index="${index}">Edit</button>
                            <button class="delete-btn" data-index="${index}">Delete</button>
                        </div>
                    </td>
                </tr>
            `;
            tbody.innerHTML += row;
        });
    
        table.appendChild(tbody);
        savedVariablesContainer.appendChild(table);
    
        // Add event listeners for delete buttons
        const deleteButtons = savedVariablesContainer.querySelectorAll('.delete-btn');
        deleteButtons.forEach(button => {
            button.addEventListener('click', function () {
                const index = parseInt(button.getAttribute('data-index'), 10);
                deleteVariable(index);
            });
        });
    
        // Add event listeners for edit buttons (optional, if you plan to implement editing)
        const editButtons = savedVariablesContainer.querySelectorAll('.edit-btn');
        editButtons.forEach(button => {
            button.addEventListener('click', function () {
                const index = parseInt(button.getAttribute('data-index'), 10);
                editVariable(index);
            });
        });
    }

    function editVariable(index) {
        // Retrieve variables from localStorage
        const savedVariables = JSON.parse(localStorage.getItem("terraformVariables")) || [];
    
        // Get the variable to edit
        const variable = savedVariables[index];
    
        // Pre-fill the form fields with the variable's data
        // (Assuming you have a form for adding/editing variables)
        const keyField = document.querySelector('input[name="key"]');
        const valueField = document.querySelector('input[name="value"]');
        const hclField = document.querySelector('input[name="hcl"]');
        const sensitiveField = document.querySelector('input[name="sensitive"]');
        const descField = document.querySelector('textarea[name="description"]');
    
        keyField.value = variable.key;
        valueField.value = variable.value;
        hclField.checked = variable.hcl;
        sensitiveField.checked = variable.sensitive;
        descField.value = variable.description || '';
    
        // Save the index being edited for updating later
        keyField.dataset.editIndex = index;
    
        // Scroll to the form (optional)
        document.querySelector('#variables-container').scrollIntoView({ behavior: 'smooth' });
    }
    
    // Handle the Save button click for edited variables
    saveButton.addEventListener('click', function () {
        const editIndex = parseInt(document.querySelector('input[name="key"]').dataset.editIndex, 10);
    
        if (!isNaN(editIndex)) {
            // Update the existing variable
            const savedVariables = JSON.parse(localStorage.getItem("terraformVariables")) || [];
            const updatedVariable = {
                key: keyField.value.trim(),
                value: valueField.value.trim(),
                hcl: hclField.checked,
                sensitive: sensitiveField.checked,
                description: descField.value.trim(),
            };
    
            savedVariables[editIndex] = updatedVariable;
            localStorage.setItem("terraformVariables", JSON.stringify(savedVariables));
    
            // Refresh the display and clear the form
            displaySavedVariables(savedVariables);
            keyField.dataset.editIndex = '';
            keyField.value = '';
            valueField.value = '';
            hclField.checked = false;
            sensitiveField.checked = false;
            descField.value = '';
        }
    });    
    
    function deleteVariable(index) {
        const confirmed = confirm("Are you sure you want to delete this variable?");
        if (!confirmed) return;
    
        // Retrieve variables from localStorage
        const savedVariables = JSON.parse(localStorage.getItem("terraformVariables")) || [];
    
        // Remove the variable at the specified index
        savedVariables.splice(index, 1);
    
        // Save the updated list back to localStorage
        localStorage.setItem("terraformVariables", JSON.stringify(savedVariables));
    
        // Refresh the display
        displaySavedVariables(savedVariables);
    
        alert("Variable has been deleted.");
    }    
    
    // Function to escape HTML to prevent XSS
    function escapeHtml(text) {
        const map = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#039;'
        };
        return text.replace(/[&<>"']/g, function(m) { return map[m]; });
    }


    // Event listener for adding a new variable
    addVariableBtn.addEventListener("click", function () {
        // Create and append a new variable block
        const newVariable = createVariableBlock();
        variablesContainer.appendChild(newVariable);

        // Show the Save button if it's the first variable added
        if (variablesContainer.children.length === 1) {
            saveButton.classList.remove("hidden");
        }
    });

    // Event listener for saving variables
    saveButton.addEventListener('click', function () {
        // Retrieve variables from localStorage
        const savedVariables = JSON.parse(localStorage.getItem("terraformVariables")) || [];
    
        // Gather all variable-block elements
        const blocks = document.querySelectorAll(".variable-block");
        const newVariables = [];
        let allValid = true;
    
        blocks.forEach((block, index) => {
            // Grab the inputs from each block
            const keyField = block.querySelector('input[name="key"]');
            const valueField = block.querySelector('input[name="value"]');
            const hclField = block.querySelector('input[name="hcl"]');
            const sensitiveField = block.querySelector('input[name="sensitive"]');
            const descField = block.querySelector('textarea[name="description"]');
    
            // Basic validation: Ensure Key and Value are not empty
            if (keyField.value.trim() === "" || valueField.value.trim() === "") {
                alert(`Please fill out both Key and Value for variable #${index + 1}.`);
                allValid = false;
                return;
            }
    
            // Build an object with each variable’s data
            const variableData = {
                key: keyField.value.trim(),
                value: valueField.value.trim(),
                hcl: hclField.checked,
                sensitive: sensitiveField.checked,
                description: descField.value.trim(),
            };
    
            newVariables.push(variableData);
        });
    
        if (!allValid) {
            return; // Stop saving if validation fails
        }
    
        // Combine existing variables with new ones
        const updatedVariables = [...savedVariables, ...newVariables];
    
        // Save the updated list back to localStorage
        localStorage.setItem("terraformVariables", JSON.stringify(updatedVariables));
    
        // Display the updated variables and clear the input fields
        displaySavedVariables(updatedVariables);
        variablesContainer.innerHTML = '';
        saveButton.classList.add('hidden');
    
        alert("Variables have been saved and displayed below.");
    
        // Scroll to the Saved Variables section
        document.querySelector('.saved-variables-section').scrollIntoView({ behavior: 'smooth' });
    });

    // Initial load of saved variables
    loadSavedVariables();
});
