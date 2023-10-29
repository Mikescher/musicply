
const API_LEVEL = {{ .APILevel }};

function playlist_iterate(obj, fn) {
    fn(obj);
    if (obj.children) for (const c of obj.children) playlist_iterate(c, fn);
}

const sleep = (milliseconds) => {return new Promise(resolve => setTimeout(resolve, milliseconds))}

function formatDuration(v) {
    return Math.floor(v/60) + ':' + `${Math.floor(v%60)}`.padStart(2, '0');
}

function formatBitrate(v) {
    return ('@ ' + Math.round(v/1000) + ' kbit');
}

//----------------------------------------------------------------------------------------------------------------------

ko.options.deferUpdates = true;

let vm = {};

//----------------------------------------------------------------------------------------------------------------------

async function loadTracks(ids) {
    try {
        vm.tracksLoading(true);
        vm.tracks([]);

        const tracks = (await Promise.all([
            await (async () => {
                let trackarr = []
                for (const plid of ids) {
                    const resp = await fetch(`/api/v${API_LEVEL}/playlists/${plid}/tracks`);
                    const rjsn = await resp.json();
                    trackarr.push(...rjsn.tracks);
                }
                return trackarr;
            })(),
            sleep(300),
        ]))[0];

        vm.tracks(tracks);

    } finally {
        vm.tracksLoading(false);
    }
}

//----------------------------------------------------------------------------------------------------------------------

vm['playlists_root'] =  {{ listPlaylists | json_indent }};
vm.playlists_root.children.unshift({ id: null, name: 'Unsorted', children: null, hasChildren: false, trackCount: 0 });
playlist_iterate(vm['playlists_root'], (obj) => { obj['active'] = ko.observable(false); });
vm.playlists_root.children[0].active(true);

vm['onPlaylistClick'] = function (pl) {
    playlist_iterate(vm['playlists_root'], (obj) => { obj.active(false); });
    pl.active(true);

    vm.tracksInitial(false);

    let ids = [];
    playlist_iterate(pl, v => ids.push(v.id));
    ids = ids.filter(p => p !== null && p !== undefined);

    loadTracks(ids).then();
}

vm['tracksLoading'] = ko.observable(false);

vm['tracksInitial'] = ko.observable(true);

vm['tracks'] = ko.observableArray();

vm['apiLevel'] = API_LEVEL;

//----------------------------------------------------------------------------------------------------------------------

ko.applyBindings(vm);