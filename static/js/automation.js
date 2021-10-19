$.ajaxSetup({
    crossDomain: false,
    beforeSend: function (xhr, settings) {
        xhr.setRequestHeader("X-CSRF-Token", $("#csrf_token").val());
    }
});

$('#gitBranchButton').on('click', function (event) {
    event.preventDefault();
    event.stopPropagation();
    let git_username = $("#usernameField").val();
    let git_password = $("#passwordField").val();
    let git_url = $("#gitUrlField").val();

    $("#gitBranchLoading").show();
    $("#gitBranchLoaded").hide();

    //fetch request
    if (git_username && git_password && git_url) {
        $.ajax({
            type: "post",
            url: "/admin/git/fetch",
            data: {
                "git_username": git_username,
                "git_password": git_password,
                "git_url": git_url,
            },
            success: function (data) {
                var gitBranchDropDown = $("#gitBranchDropDown");
                $.each(data["branches"], function (index, branch) {
                    gitBranchDropDown.append(new Option(branch, branch));
                });
            },
            complete: function(){
                $('#gitBranchLoaded').show();
                $('#gitBranchLoading').hide();
            }
        });
    }

});