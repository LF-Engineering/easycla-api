class PostgresSchema(object):
    """
    Class that handles mapping dynamo db instances to the pg tables
    """
    def __init__(self, model):
        self.model = model
        self.table_mapping = {
            "users": self.users,
            "gerrit_instances": self.gerrit_instances,
            "companies": self.companies,
            "repositories": self.repositories,
            "github_orgs": self.github_orgs,
            "projects": self.projects,
            "signatures": self.signatures,
        }
        self.committed = False

    def get_tables(self, table):
        return self.table_mapping[table]

    def set_committed(self, committed):
        self.committed = committed

    def get_committed(self):
        return self.committed

    def users(self):
        return {
            "user_id": self.model.model.user_id,
            "lf_email": self.model.model.lf_email,
            "lf_username": self.model.model.lf_username,
            "user_github_id": self.model.model.user_github_id,
            "user_company_id": self.model.model.user_company_id,
            "user_name": self.model.model.user_name,
            "user_github_name": self.model.model.user_github_username,
        }

    def gerrit_instances(self):
        return {
            "gerrit_id": self.model.model.gerrit_id,
            "date_created": self.model.model.date_created,
            "date_modified": self.model.model.date_modified,
            "gerrit_name": self.model.model.gerrit_name,
            "gerrit_url": self.model.model.gerrit_url,
            "group_id_ccla": self.model.model.group_id_ccla,
            "group_id_icla": self.model.model.group_id_icla,
            "group_name_ccla": self.model.model.group_name_ccla,
            "group_name_icla": self.model.model.group_name_icla,
            "project_id": self.model.model.project_id,
            "version": self.model.model.version,
        }

    def companies(self):
        return {
            "company_id": self.model.model.company_id,
            "company_manager_id": self.model.model.company_manager_id,
            "company_name": self.model.model.company_name,
        }

    def projects(self):
        return {
            "project_id": self.model.model.project_id,
            # "project_acl": model.model.project_acl,
            "project_ccla_enabled": self.model.model.project_ccla_enabled,
            "project_ccla_requires_icla_signature": self.model.model.project_ccla_requires_icla_signature,
            # "project_corporate_documents": model.model.project_corporate_documents,
            "project_icla_enabled": self.model.model.project_icla_enabled,
            # "project_individual_documents": model.model.project_individual_documents,
            "project_name": self.model.model.project_name,
        }

    def signatures(self):
        return {
            "signature_id": self.model.model.signature_id,
            "signature_approved": self.model.model.signature_approved,
            "signature_document_major_version": self.model.model.signature_document_major_version,
            "signature_document_minor_version": self.model.model.signature_document_minor_version,
            "signature_project_id": self.model.model.signature_project_id,
            "signature_reference_id": self.model.model.signature_reference_id,
            "signature_reference_type": self.model.model.signature_reference_type,
            "signature_signed": self.model.model.signature_signed,
            "signature_type": self.model.model.signature_type,
            "signature_user_ccla_company_id": self.model.model.signature_user_ccla_company_id,
            "signature_acl": self.model.model.signature_acl,
            "domain_whitelist": self.model.model.domain_whitelist,
            "signature_return_url": self.model.model.signature_return_url,
            # "email_whitelist": model.model.email_whitelist,
            # "github_whitelist": model.model.github_whitelist,
            # "github_org_whitelist": model.model.github_org_whitelist,
        }

    def repositories(self):
        return {
            "repository_id": self.model.model.repository_id,
            "repository_external_id": self.model.model.repository_external_id,
            "repository_organisation_name": self.model.model.repository_organization_name,
            "repository_project_id": self.model.model.repository_project_id,
            "repository_sfdc_id": self.model.model.repository_sfdc_id,
            "repository_type": self.model.model.repository_type,
            "repository_url": self.model.model.repository_url,
            "version": self.model.model.version,
        }

    def github_orgs(self):
        return {
            "organization_name": self.model.model.organization_name,
            "organization_sfid": self.model.model.organization_sfid,
            "version": self.model.model.version,
            "organization_installation_id": self.model.model.organization_installation_id,
        }
