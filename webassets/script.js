
/*{{ "const API_LEVEL =" | safe }}*/ /*{{ .APILevel }}*/;

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
    while (currentIndex > 0) {
        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex--;
        [array[currentIndex], array[randomIndex]] = [array[randomIndex], array[currentIndex]];
    }
    return array;
}

//----------------------------------------------------------------------------------------------------------------------

ko.options.deferUpdates = true;

let vm = {};

//----------------------------------------------------------------------------------------------------------------------

const aplayer = new Audio();

aplayer.addEventListener('error', function (evt) {
    if (vm.playbackStatus() === 'playing') vm.playbackStatus('paused');
    vm.playbackProgress(0);
    vm.playbackTotal(null);
    console.error('aplayer::error', evt);
});

aplayer.addEventListener('loadeddata', function (evt) {
    vm.playbackProgress(aplayer.currentTime);
    vm.playbackTotal(aplayer.duration);
});

aplayer.addEventListener('timeupdate', function (evt) {
    vm.playbackProgress(aplayer.currentTime);
});

aplayer.addEventListener('durationchange', function (evt) {
    vm.playbackTotal(aplayer.duration);
});

aplayer.addEventListener('seeked', function (evt) {
    vm.playbackProgress(aplayer.currentTime);
});

aplayer.addEventListener('ended', function (evt) {
    playNext().then();
});

aplayer.addEventListener('play', function (evt) {
    vm.playbackStatus('playing');
});

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
        await changeActiveQueueItem(id);

    } else if (vm.queue().at(-1).type() === 'past') {

        vm.queue.push({
            queueID: id,
            type: ko.observable('active'),
            track: track,
        });
        await changeActiveQueueItem(id);

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

    if (vm.playbackStatus() === 'finished') vm.playbackStatus('paused');

    return id;
}

async function changeActiveQueueItem(trackid) {

    aplayer.pause();
    vm.playbackStatus('paused');

    let past = true;

    for (const e of vm.queue()) {

        if (e.queueID === trackid) {

            e.type('active');
            vm.playbackTotal(null);
            vm.playbackProgress(0);
            vm.playbackStatus('paused');

            await scrollQueueEntryIntoView(e.queueID);

            aplayer.src = `/api/v${API_LEVEL}/playlists/${e.track.playlistID}/tracks/${e.track.id}/stream`;

            past = false;

        } else if (past) {

            e.type('past');

        } else {

            e.type('future');

        }
    }
}

async function playNext() {
    let found = false;
    for (const e of vm.queue()) {
        if (found) {
            await changeActiveQueueItem(e.queueID);
            await aplayer.play();
            return;
        } else if (e.type() === 'active') {
            found = true;
        }
    }

    // else (last queue item)

    for (const e of vm.queue()) e.type('past');

    aplayer.pause();
    vm.playbackStatus('finished');
    vm.playbackProgress(0);
    vm.playbackTotal(null)
}

async function scrollQueueEntryIntoView(queueid) {
    document.querySelector(`.queue_item[data-queue-entry-id="${queueid}"]`)?.scrollIntoView({behavior: 'smooth', block: 'nearest'});
}

//----------------------------------------------------------------------------------------------------------------------

/*{{ "vm['playlists_root'] =" | safe }}*/ /*{{ listPlaylists | json_indent }}*/;

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
    (async () =>
    {
        aplayer.pause()
        vm.playbackStatus('paused');
        for (const t of vm.tracks()) await enqueue(t);
        await aplayer.play();
    })().then();
};

vm['onShuffle'] = function () {
    (async () =>
    {
        aplayer.pause()
        vm.playbackStatus('paused');
        for (const t of shuffle(vm.tracks())) await enqueue(t);
        await aplayer.play();
    })().then();
};

vm['onPlaySingle'] = function (track) {
    (async () =>
    {
        vm.queue.removeAll();
        vm.playbackStatus('finished');
        vm.playbackProgress(0);

        aplayer.pause();
        vm.playbackStatus('paused');

        await enqueue(track);
        await aplayer.play();
    })().then();
};

vm['onEnqueueSingle'] = function (track) {
    enqueue(track).then();
};

vm['onQueueClear'] = function () {
    vm.queue.removeAll();

    aplayer.pause();
    vm.playbackStatus('finished');
    vm.playbackProgress(0);
    vm.playbackTotal(null);
}

vm['onPlaybackPlay'] = function () {
    aplayer.play().then();
}

vm['onPlaybackPause'] = function () {
    aplayer.pause();
    vm.playbackStatus('paused');
}

vm['onPlaybackRestart'] = function () {
    (async () =>
    {
        aplayer.pause()
        vm.playbackStatus('paused');
        await changeActiveQueueItem(vm.queue()[0].queueID);
        await aplayer.play();
    })().then();
}

vm['onPlayPastTrack'] = function (queueitem) {
    (async () =>
    {
        aplayer.pause()
        vm.playbackStatus('paused');
        await changeActiveQueueItem(queueitem.queueID);
        await aplayer.play();
    })().then();
}

vm['onPlayFutureTrack'] = function (queueitem) {
    (async () => {
        aplayer.pause()
        vm.playbackStatus('paused');
        await changeActiveQueueItem(queueitem.queueID);
        await aplayer.play();
    })().then();
}

vm['onPlayNextTrack'] = function () {
    playNext().then();
}

vm['onManualSeek'] = function (_, evt) {
    const tseek = evt.target.valueAsNumber
    vm.playbackProgress(tseek);
    aplayer.currentTime = tseek;
}

//----------------------------------------------------------------------------------------------------------------------

ko.applyBindings(vm);