import { Form, FormikHelpers, FormikProvider } from 'formik';

import Button from '@/components/ui/Button';
import Select from '@/components/ui/Select';
import TextField from '@/components/ui/TextField';
import { ROLE_OPTIONS, useNewMemberForm } from '@/features/members/types/addMemberFormService';
import type { NewMemberDTO } from '@/features/members/types/memberServices';

type NewMemberFormProps = {
  onSubmit: (values: NewMemberDTO, helpers: FormikHelpers<NewMemberDTO>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
};

function AddNewMemberForm({ onSubmit, isSubmitting }: NewMemberFormProps) {
  const formik = useNewMemberForm(onSubmit);

  const { touched, errors, getFieldProps, setFieldValue, values } = formik;
  return (
    <FormikProvider value={formik}>
      <Form className="space-y-3">
        <TextField
          label="Email Address"
          placeholder="member_email@domain.com"
          errorMessage={touched.email && errors.email}
          {...getFieldProps('email')}
        />
        <fieldset>
          <label htmlFor="role" className="text-sm">
            Select role
          </label>
          <Select
            id="role"
            isDisabled={isSubmitting}
            defaultValue={ROLE_OPTIONS.filter((opt) => opt.value === values.role)}
            options={ROLE_OPTIONS.filter((opt) => opt.value !== values.role)}
            name="role"
            value={ROLE_OPTIONS.filter((opt) => opt.value === values.role)}
            onChange={(value: any) => setFieldValue('role', value.value)}
          />
        </fieldset>
        <div className="pt-3 text-center">
          <Button type="submit" isLoading={isSubmitting} isDisabled={isSubmitting}>
            Invite
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
}

export default AddNewMemberForm;
