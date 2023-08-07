$(document).ready(function() {
    $("#submit").click(function(e) {
        e.preventDefault();
        var link = $("#link").val();
        $.ajax({
            url: "http://ama1.ru/shortener/",
            type: "POST",
            data: link,
            success: function(response) {
                $("#shortened-link").text(response);
                $("#shortened-link-alert").text(response ? "" : "An error has occurred, check the correctness of the link.");
            },
            error: function(xhr, status, error) {
                $("#shortened-link").text("");
                $("#shortened-link-alert").text(error || "An error has occurred, check the correctness of the link.");
            }
        });
    });
});
$(document).ready(function() {
    $("#shortened-link").click(function() {
        var linkText = $(this).text();
        navigator.clipboard.writeText(linkText)
            .then(function() {
                alert("Link copied to clipboard!");
            })
            .catch(function(error) {
                console.error("Error copying link to clipboard: ", error);
            });
    });
});