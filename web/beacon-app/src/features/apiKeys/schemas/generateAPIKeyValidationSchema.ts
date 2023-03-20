import * as Yup from 'yup';

const generateAPIKeyValidationSchema = Yup.object().shape({
  name: Yup.string().required('The key name is required.'),
});

export default generateAPIKeyValidationSchema;
