import xhook from 'xhook';
import $ from 'jquery';
window.jQuery = $;
window.$ = $;
import 'popper.js';
import 'parsleyjs';
import 'bootstrap';
//import 'jquery-slimscroll';
import 'select2';

import '../scss/application.scss'

import fontawesome from "@fortawesome/fontawesome";
import fasCog from "@fortawesome/fontawesome-free-solid/faCog";
import fasPencil from "@fortawesome/fontawesome-free-solid/faPencilAlt";
import fasTags from "@fortawesome/fontawesome-free-solid/faTags";
import fasCopy from "@fortawesome/fontawesome-free-solid/faCopy";
import fasUser from "@fortawesome/fontawesome-free-solid/faUser";
import fasPlay from "@fortawesome/fontawesome-free-solid/faPlay";
import fasSignOut from "@fortawesome/fontawesome-free-solid/faSignOutAlt";
import fasEdit from "@fortawesome/fontawesome-free-solid/faEdit";
import fasTimes from "@fortawesome/fontawesome-free-solid/faTimes";
import fasSignIn from "@fortawesome/fontawesome-free-solid/faSignInAlt";
import fasUserPlus from "@fortawesome/fontawesome-free-solid/faUserPlus";
import fasArchive from "@fortawesome/fontawesome-free-solid/faArchive";
import fasHome from "@fortawesome/fontawesome-free-solid/faHome";
import fasEye from "@fortawesome/fontawesome-free-solid/faEye";
import fasCheck from "@fortawesome/fontawesome-free-solid/faCheck";
import fasCalendar from "@fortawesome/fontawesome-free-solid/faCalendarAlt";
import fasChevronRight from "@fortawesome/fontawesome-free-solid/faChevronRight";
import fasComments from "@fortawesome/fontawesome-free-solid/faComments";
import fasExclamationTriangle from "@fortawesome/fontawesome-free-solid/faExclamationTriangle";

fontawesome.library.add(fasCog, fasPencil, fasTags, fasComments, fasExclamationTriangle, fasCopy, fasCheck, fasChevronRight, fasCalendar, fasEye, fasUser, fasPlay, fasSignOut, fasEdit, fasTimes, fasSignIn, fasUserPlus, fasArchive, fasHome);

$(document).ready(function () {
    $('select#tags').select2({
        tags: true,
        tokenSeparators: [','],
    });

    if (document.querySelector('#ck-content')) {
        //add csrf protection to ckeditor uploads
        xhook.before(function (request) {
            if (!/^(GET|HEAD|OPTIONS|TRACE)$/i.test(request.method)) {
                request.xhr.setRequestHeader("X-CSRF-TOKEN", window.csrf_token);
            }
        });

        ClassicEditor
            .create(document.querySelector('#ck-content'), {
                language: 'en', //to set different lang include <script src="/public/js/ckeditor/build/translations/{lang}.js"></script> along with core ckeditor script
                ckfinder: {
                    uploadUrl: '/admin/upload'
                },
            })
            .catch(error => {
                console.error(error);
            });
    }
});