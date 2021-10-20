const addNewFileButton = document.getElementById("addNewFile")
if (addNewFileButton) {
        addNewFileButton.addEventListener('click', () => {

        // main container
        let fileContainer = document.getElementById('filesContainer');
        fileCount = fileContainer.children.length + 1

        console.log('fileContainer--->', fileCount);
        let headingWrappingDiv = document.createElement('div');
        headingWrappingDiv.setAttribute('class', 'my-2')
        let fileHeading = document.createElement('b');
        fileHeading.innerText = 'File ' + fileCount;
        headingWrappingDiv.appendChild(fileHeading);

        let wrappingDiv = document.createElement('div');
        wrappingDiv.setAttribute('class', 'p-3 rounded-3 border border-light bg-light mb-3')

        // creating filename container
        let fileName = document.createElement('input');
        fileName.setAttribute('class', 'form-control');
        fileName.setAttribute('type', 'text');
        fileName.setAttribute('name', 'fileNameField' + fileCount);
        fileName.setAttribute('id', 'fileNameField' + fileCount);
        fileName.setAttribute('placeholder', 'File Name');
        let fileNameLabel = document.createElement('label');
        fileNameLabel.setAttribute('for', 'fileNameField' + fileCount);
        fileNameLabel.innerText = "File Name";
        let fileNameContainerDiv = document.createElement('div');
        fileNameContainerDiv.setAttribute('class', 'login-input-primary form-floating mb-3 pe-1 col');
        fileNameContainerDiv.appendChild(fileName);
        fileNameContainerDiv.appendChild(fileNameLabel);

        // FileTemplate
        let fileTemplate = document.createElement('textarea');
        fileTemplate.setAttribute('class', 'form-control')
        fileTemplate.setAttribute('type', 'text');
        fileTemplate.setAttribute('name', 'fileTemplateField' + fileCount)
        fileTemplate.setAttribute('id', 'fileTemplateField' + fileCount);
        fileTemplate.setAttribute('placeholder', 'File Template');
        fileTemplate.setAttribute('style', 'height: 200px')
        let fileTemplateLabel = document.createElement('label');
        fileTemplateLabel.setAttribute('for', 'fileTemplateField' + fileCount);
        fileTemplateLabel.innerText = "File Template";
        let fileTemplateContainerDiv = document.createElement('div');
        fileTemplateContainerDiv.setAttribute('class', 'login-input-primary form-floating mb-3 pe-1 col');
        fileTemplateContainerDiv.appendChild(fileTemplate);
        fileTemplateContainerDiv.appendChild(fileTemplateLabel);

        // TemplateMetadata
        let fileMapping = document.createElement('textarea');
        fileMapping.setAttribute('class', 'form-control')
        fileMapping.setAttribute('type', 'text');
        fileMapping.setAttribute('name', 'fileMappingField' + fileCount)
        fileMapping.setAttribute('id', 'fileMappingField' + fileCount);
        fileMapping.setAttribute('placeholder', 'File Metadata');
        fileMapping.setAttribute('style', 'height: 200px')
        let fileMappingLabel = document.createElement('label');
        fileMappingLabel.setAttribute('for', 'fileMappingField' + fileCount);
        fileMappingLabel.innerText = "File Metadata";
        let fileMappingContainerDiv = document.createElement('div');
        fileMappingContainerDiv.setAttribute('class', 'login-input-primary form-floating mb-3 pe-1 col');
        fileMappingContainerDiv.appendChild(fileMapping);
        fileMappingContainerDiv.appendChild(fileMappingLabel);

        wrappingDiv.appendChild(fileHeading);
        wrappingDiv.appendChild(fileNameContainerDiv);
        wrappingDiv.appendChild(fileTemplateContainerDiv);
        wrappingDiv.appendChild(fileMappingContainerDiv);

        fileContainer.appendChild(wrappingDiv);
    });
}
