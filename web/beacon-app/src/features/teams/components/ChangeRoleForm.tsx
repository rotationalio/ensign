import { t, Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { Form, Formik, FormikHelpers } from 'formik';

import Select from '@/components/ui/Select';

import { ChangeRoleFormDto } from '../types/changeRoleFormDto';

const ROLE_OPTIONS = [
  { value: 'Owner', label: 'Owner' },
  { value: 'Admin', label: 'Admin' },
  { value: 'Member', label: 'Member' },
  { value: 'Observer', label: 'Observer' },
];

type ChangeRoleFormProps = {
  handleSubmit: (values: ChangeRoleFormDto, helpers: FormikHelpers<ChangeRoleFormDto>) => void;
  initialValues: ChangeRoleFormDto;
};

const ChangeRoleForm = ({ handleSubmit, initialValues }: ChangeRoleFormProps) => {
  return (
    <Formik onSubmit={handleSubmit} initialValues={initialValues} enableReinitialize>
      {({ getFieldProps, values, setFieldValue, isSubmitting }: any) => (
        <Form className="space-y-3">
          <TextField
            label={t`Team Member`}
            placeholder="Natali Craig"
            {...getFieldProps('name')}
            isDisabled
            data-cy="teamMemberName"
            fullWidth
          />
          <TextField
            label={t`Current role`}
            placeholder="Member"
            {...getFieldProps('current_role')}
            isDisabled
            fullWidth
            data-cy="teamMemberRole"
          />
          <fieldset>
            <label htmlFor="role" className="text-sm">
              <Trans>Select new role</Trans>
            </label>
            <Select
              id="role"
              isDisabled={isSubmitting}
              defaultValue={ROLE_OPTIONS.filter((opt) => opt.value === values.current_role)}
              options={ROLE_OPTIONS.filter((opt) => opt.value !== values.current_role)}
              name="role"
              value={ROLE_OPTIONS.filter((opt) => opt.value === values.role)}
              onChange={(value: any) => setFieldValue('role', value.value)}
            />
          </fieldset>
          <div className="pt-3 text-center">
            <Button
              type="submit"
              isLoading={isSubmitting}
              disabled={isSubmitting}
              data-cy="saveNewRole"
            >
              <Trans>Save</Trans>
            </Button>
          </div>
        </Form>
      )}
    </Formik>
  );
};

export default ChangeRoleForm;
