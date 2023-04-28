import { t } from '@lingui/macro';
import * as Yup from 'yup';

const createProjectTopicSchema = Yup.object().shape({
  topic_name: Yup.string()
    .trim()
    .required(t`Topic name is required.`)
    .matches(/^[^\s]*$/, t`Topic name cannot include spaces.`)
    .matches(/^[^_-].*$/, t`Topic name cannot start with an underscore or dash.`)
    .max(512, t`Topic name must be less than 512 characters.`),
});

export default createProjectTopicSchema;
