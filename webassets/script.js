
ko.options.deferUpdates = true;

let vm = {};

vm['playlists_root'] = {{ listPlaylists | json_indent }};

ko.applyBindings(vm);