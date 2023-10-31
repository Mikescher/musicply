
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

function shuffle(array) {
    let currentIndex = array.length,  randomIndex;

    // While there remain elements to shuffle.
    while (currentIndex > 0) {

        // Pick a remaining element.
        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex--;

        // And swap it with the current element.
        [array[currentIndex], array[randomIndex]] = [
            array[randomIndex], array[currentIndex]];
    }

    return array;
}

//----------------------------------------------------------------------------------------------------------------------

ko.options.deferUpdates = true;

let vm = {};

//----------------------------------------------------------------------------------------------------------------------

const aplayer = new Audio();

aplayer.addEventListener('error', function (evt) {
    //TODO
})

aplayer.addEventListener('abort', function (evt) {
    //TODO
})

aplayer.addEventListener('loadeddata', function (evt) {
    vm.playbackProgress(aplayer.currentTime);
    vm.playbackTotal(aplayer.duration);

})

aplayer.addEventListener('timeupdate', function (evt) {
    vm.playbackProgress(aplayer.currentTime);
})

aplayer.addEventListener('durationchange', function (evt) {
    vm.playbackTotal(aplayer.duration);
})

aplayer.addEventListener('seeked', function (evt) {
    vm.playbackProgress(aplayer.currentTime);
})

aplayer.addEventListener('ended', function (evt) {
    //TODO
})

aplayer.addEventListener('pause', function (evt) {
    vm.playbackStatus('paused');
})

aplayer.addEventListener('play', function (evt) {
    vm.playbackStatus('playing');
})

//----------------------------------------------------------------------------------------------------------------------

const uniqueid = () => Math.ceil(Math.random() * 1000000000).toString(16).toUpperCase().padStart(8, '0');

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

async function searchTracks(q) {
    try {
        vm.tracksLoading(true);
        vm.tracks([]);

        const tracks = (await Promise.all([
            await (async () => {
                if (q.trim() === '') return [];
                const resp = await fetch(`/api/v${API_LEVEL}/tracks?search=${encodeURI(q)}`);
                const rjsn = await resp.json();
                return [...rjsn.tracks];
            })(),
            sleep(300),
        ]))[0];

        vm.tracks(tracks);

    } finally {
        vm.tracksLoading(false);
    }
}

async function searchPlaylistTracks(ids, q) {
    try {
        vm.tracksLoading(true);
        vm.tracks([]);

        const tracks = (await Promise.all([
            await (async () => {
                let trackarr = []
                for (const plid of ids) {
                    const resp = await fetch(`/api/v${API_LEVEL}/playlists/${plid}/tracks?search=${encodeURI(q)}`);
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

async function enqueue(track) {

    const id = uniqueid();

    if (vm.queue().length === 0) {

        vm.queue.push({
            queueID: id,
            type: ko.observable('active'),
            track: track,
        });
        await changeActiveTrack(id);

    } else if (vm.queue().at(-1).type() === 'past') {

        vm.queue.push({
            queueID: id,
            type: ko.observable('active'),
            track: track,
        });
        await changeActiveTrack(id);

    } else if (vm.queue().at(-1).type() === 'active') {

        vm.queue.push({
            queueID: id,
            type: ko.observable('future'),
            track: track,
        });

    } else if (vm.queue().at(-1).type() === 'future') {

        vm.queue.push({
            queueID: id,
            type: ko.observable('future'),
            track: track,
        });

    }

    if (vm.playbackStatus() === 'finished') {
        vm.playbackStatus('paused');

        aplayer.pause();
    }
}

async function changeActiveTrack(trackid) {

    aplayer.pause();
    aplayer.src = '';

    let past = true;

    for (const e of vm.queue()) {

        if (e.queueID === trackid) {

            e.type('active');
            vm.playbackTotal(null);
            vm.playbackProgress(0);
            vm.playbackStatus('paused')

            aplayer.src = `/api/v${API_LEVEL}/playlists/${e.track.playlistID}/tracks/${e.track.id}/stream`;

            past = false;

        } else if (past) {

            e.type('past');

        } else {

            e.type('future');

        }
    }
}

//----------------------------------------------------------------------------------------------------------------------

vm['playlists_root'] =  {{ listPlaylists | json_indent }};
vm.playlists_root.children.unshift({ id: null, name: 'All', children: null, hasChildren: false, trackCount: 0 });
playlist_iterate(vm['playlists_root'], (obj) => { obj['active'] = ko.observable(false); });

vm['tracksLoading'] = ko.observable(false);

vm['tracksInitial'] = ko.observable(true);

vm['tracks'] = ko.observableArray(); // Track[]

vm['queue'] = ko.observableArray(); // { queueID: string, type: ('past'|'active'|'future'), track: Track }[]

vm['apiLevel'] = API_LEVEL;

vm['searchText'] = ko.observable();

vm['playbackStatus'] = ko.observable('finished'); // 'playing'|'paused'|'finished'

vm['playbackProgress'] = ko.observable(0);

vm['playbackTotal'] = ko.observable(null);

//----------------------------------------------------------------------------------------------------------------------

vm['onSearchKeyPress'] = function (data, event) {
    if (event.keyCode === 13) vm.onSearch();
    return true;
};

vm['onPlaylistClick'] = function (pl) {
    playlist_iterate(vm['playlists_root'], (obj) => { obj.active(false); });
    pl.active(true);

    vm.searchText('');

    vm.tracksInitial(false);

    let ids = [];
    playlist_iterate(pl, v => ids.push(v.id));
    ids = ids.filter(p => p !== null && p !== undefined);

    loadTracks(ids).then();
}

vm['onSearch'] = function () {
    vm.tracksInitial(false);

    let ids = [];
    playlist_iterate(vm.playlists_root, plst => {
        if (plst.active()) playlist_iterate(plst, v => ids.push(v.id));
    });

    if (vm.playlists_root.children[0].active()) {
        searchTracks(vm.searchText()).then();
    } else if (ids.length === 0) {
        vm.playlists_root.children[0].active(true);
        searchTracks(vm.searchText()).then();
    } else {
        ids = ids.filter(p => p !== null && p !== undefined);
        searchPlaylistTracks(ids, vm.searchText()).then();
    }
};

vm['onPlayAll'] = function () {
    for (const t of vm.tracks()) enqueue(t).then();
};

vm['onShuffle'] = function () {
    for (const t of shuffle(vm.tracks())) enqueue(t).then();
};

vm['onPlaySingle'] = function (track) {
    //TODO
};

vm['onEnqueueSingle'] = function (track) {
    enqueue(track).then();
};

vm['onQueueClear'] = function () {
    vm.queue.removeAll();
    //TODO
    vm.playbackStatus('finished');
    vm.playbackProgress(0);

    aplayer.pause();
    aplayer.src = '';
}

vm['playbackPlay'] = function () {
    aplayer.play().then();
}

vm['playbackPause'] = function () {
    aplayer.pause();
}

vm['onPlayPastTrack'] = function () {
    //TODO
}

vm['onPlayFutureTrack'] = function () {
    //TODO
}

vm['onManualSeek'] = function (_, evt) {
    const tseek = evt.target.valueAsNumber
    vm.playbackProgress(tseek);
    aplayer.currentTime = tseek;
}

//----------------------------------------------------------------------------------------------------------------------

ko.applyBindings(vm);