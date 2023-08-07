$(document).ready(function() {
    $("#submit").click(function(e) {
        e.preventDefault();
        var link = $("#link").val();
        $.ajax({
            url: "http://127.0.0.1:8083/shortener/",
            type: "POST",
            data: link,
            success: function(response) {
                $("#shortened-link").text(response);
                $("#shortened-link-alert").text(response ? "" : "Произошла ошибка(500)");
            },
            error: function(xhr, status, error) {
                $("#shortened-link").text("");
                $("#shortened-link-alert").text(error || "Произошла ошибка, проверьте корректность ссылки.");
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