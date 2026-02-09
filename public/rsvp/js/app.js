// search invitation

var invitations = [];
var selectedInvitations = [];

const searchContainer = $(".search-container");
const selectContainer = $(".select-container");
const confirmContainer = $(".confirm-container");
const thanksContainer = $(".thanks-container");

searchContainer.show();
selectContainer.hide();
confirmContainer.hide();
thanksContainer.hide();
var path = window.location.pathname;
var evento = path.split("/")[2] || "default";


function searchInvitation() {
    var search = $(".search-container__input").val();
    if (search.length < 3) {
        return;
    }
    $.ajax({
        url: "/api/search/" + evento,
        type: "GET",
        data: { search: search },
        success: function (response) {
            invitations = response;
            showInvitations(response);
        }
    });
}

function showInvitations(invitations) {
    searchContainer.hide();
    selectContainer.show();
    var invitationsHtml = "";
    if (invitations.length == 0) {
        invitationsHtml += "<div class='select-container__select__item'>";
        invitationsHtml += "<label for='invitation-0'>No se encontraron invitaciones</label>";
        invitationsHtml += "</div>";
    }
    invitations.forEach(function (invitation) {
        invitationsHtml += "<div class='select-container__select__item'>";
        invitationsHtml += "<input type='checkbox' class='select-container__checkbox' id='invitation-" + invitation.ID + "'>";
        invitationsHtml += "<label for='invitation-" + invitation.ID + "'>" + invitation.nombre + " " + invitation.apellido + " (" + invitation.code + ")" + "</label>";
        invitationsHtml += "</div>";
    });
    $(".select-container__select").html(invitationsHtml);
}

function selectInvitations() {
    selectContainer.hide();
    confirmContainer.show();
    var invitationsChecked = $(".select-container__checkbox:checked");
    var invitationsHtml = "";
    invitationsChecked.each(function () {
        var id = this.id.replace("invitation-", "");
        var invitation = invitations.find(function (invitation) {
            return invitation.ID == id;
        });
        selectedInvitations.push(invitation);
        invitationsHtml += "<div class='confirm-container__invitations__item'>";
        invitationsHtml += "<label for='invitation-" + invitation.ID + "'>" + invitation.nombre + " " + invitation.apellido + " (" + invitation.code + ")" + " [" + invitation.respuesta + "]</label>";
        invitationsHtml += "</div>";

    });
    $(".confirm-container__invitations").html(invitationsHtml);
    $(".select-container").hide();
    $(".confirm-container").show();
}

function sendConfirmation(invitationId, respuesta) {
    $.ajax({
        url: "/api/invitados/" + invitationId,
        type: "PATCH",
        data: JSON.stringify({ Respuesta: respuesta }),
        contentType: "application/json",
        success: function (response) {
            console.log(response);
        }
    });
}

function confirmationYes() {
    selectedInvitations.forEach(function (invitation) {
        sendConfirmation(invitation.ID, "si asistira");
    });
    $(".confirm-container").hide();
    $(".thanks-container__text").text("Esperamos verte en nuestra celebración");
    $(".thanks-container").show();
    setTimeout(() => {
        $(".search-container").show();
        $(".thanks-container").hide();
    }, 5000);
}
function confirmationNo() {
    selectedInvitations.forEach(function (invitation) {
        sendConfirmation(invitation.ID, "no asistira");
    });
    $(".confirm-container").hide();
    $(".thanks-container__text").text("Lamentamos que no puedas acompañarnos en nuestra celebración");
    $(".thanks-container").show();
    setTimeout(() => {
        $(".search-container").show();
        $(".thanks-container").hide();
    }, 5000);
}

function backToSearch() {
    searchContainer.show();
    selectContainer.hide();
    confirmContainer.hide();
    thanksContainer.hide();
}

function backToSelect() {
    selectContainer.show();
    confirmContainer.hide();
    thanksContainer.hide();
}

$(".back-to-search").click(backToSearch);
$(".back-to-select").click(backToSelect);
$(".confirm-container__button--yes").click(confirmationYes);
$(".confirm-container__button--no").click(confirmationNo);
$(".search-container__button").click(searchInvitation);
$(".select-container__button").click(selectInvitations);
$(".back-to-search").click(backToSearch);
$(".back-to-select").click(backToSelect);
// ON ENTER
$(".search-container__input").on("keyup", function (e) {
    if (e.key == "Enter") {
        searchInvitation();
    }
});
