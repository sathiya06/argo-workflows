import {ObjectMeta} from 'argo-ui/src/models/kubernetes';
import {useEffect, useState} from 'react';
import * as React from 'react';

import {Link, Workflow} from '../models';
import {services} from '../services';
import {Button} from './button';

function toEpoch(datetime: string) {
    if (datetime) {
        return new Date(datetime).getTime();
    } else {
        return Date.now();
    }
}

function addEpochTimestamp(jsonObject: {metadata: ObjectMeta; workflow?: Workflow; status?: any}) {
    if (jsonObject === undefined || jsonObject.status.startedAt === undefined) {
        return;
    }

    jsonObject.status.startedAtEpoch = toEpoch(jsonObject.status.startedAt);
    jsonObject.status.finishedAtEpoch = toEpoch(jsonObject.status.finishedAt);
}

function splitWithMetadataKnowledge(replaceable: string) {
    const parts = replaceable.split('.');
    if (replaceable.startsWith('workflow.metadata.labels') || replaceable.startsWith('workflow.metadata.annotations')) {
        // Take the first 3 parts, then join the rest
        const result = parts.slice(0, 3);
        result.push(parts.slice(3).join('.'));

        return result;
    }
    return parts;
}

export function processURL(urlExpression: string, jsonObject: any) {
    addEpochTimestamp(jsonObject);
    /* replace ${} from input url with corresponding elements from object
    only return null for known variables, otherwise empty string*/
    return urlExpression.replace(/\${[^}]*}/g, x => {
        const replaced = x.replace(/(\$%7B|%7D|\${|})/g, '');
        const parts = splitWithMetadataKnowledge(replaced);
        const emptyVal = parts[0] === 'workflow' ? '' : null;
        const res = parts.reduce((p: any, c: string) => (p && p[c]) || emptyVal, jsonObject);
        return res;
    });
}

export function openLinkWithKey(url: string, target?: string) {
    if ((window.event as MouseEvent).ctrlKey || (window.event as MouseEvent).metaKey) {
        window.open(url, '_blank');
    } else if (target !== `''`) {
        window.open(url, target);
    } else {
        document.location.href = url;
    }
}

export function Links({scope, object, button}: {scope: string; object: {metadata: ObjectMeta; workflow?: Workflow; status?: any}; button?: boolean}) {
    const [links, setLinks] = useState<Link[]>();
    const [error, setError] = useState<Error>();
    useEffect(() => {
        services.info
            .getInfo()
            .then(x => (x.links || []).filter(y => y.scope === scope))
            .then(setLinks)
            .catch(setError);
    }, []);

    return (
        <>
            {error && error.message}
            {links &&
                links.map(({url, name, target}) => {
                    if (button) {
                        return (
                            <Button onClick={() => openLinkWithKey(processURL(url, object), target)} key={name} icon='external-link-alt'>
                                {name}
                            </Button>
                        );
                    }
                    return (
                        <a key={name} href={processURL(url, object)} target={target} rel='noreferrer'>
                            {name} <i className='fa fa-external-link-alt' />
                        </a>
                    );
                })}
        </>
    );
}
