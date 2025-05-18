#ifndef ATTACHMENTS_COMP_H
#define ATTACHMENTS_COMP_H

#include <QWidget>
#include <QVBoxLayout>
#include <QByteArray>
#include <QSaveFile>
#include <QStandardPaths>
#include <QMessageBox>
namespace Ui {
class attachments_comp;
}

class attachments_comp : public QWidget
{
    Q_OBJECT

public:
    explicit attachments_comp(QWidget *parent = nullptr);
    ~attachments_comp();

    void setData(const QString &name, QByteArray data);

private slots:
    void on_pushButton_clicked();

private:
    Ui::attachments_comp *ui;
    QByteArray data;
};

#endif // ATTACHMENTS_COMP_H
